package server_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestWatchdog_DetectsTampering(t *testing.T) {
	callCount := int32(0)
	var tamperDetected atomic.Bool

	fakeManager := &testutils.FakeManager{
		VerifyLicenseSeatsFn: func(ctx context.Context) error {
			atomic.AddInt32(&callCount, 1)
			return fmt.Errorf("license abc: seat count 99 exceeds maximum of 3 (possible tampering detected)")
		},
	}

	cfg := &server.Config{
		VerifyInterval: 50 * time.Millisecond,
	}

	watchdog := server.NewWatchdog(cfg, fakeManager)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		watchdog.Start(ctx, func() {
			tamperDetected.Store(true)
		})
		close(done)
	}()

	<-done

	assert.True(t, tamperDetected.Load(), "onTamper callback should have been called")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&callCount), int32(1), "VerifyLicenseSeats should have been called at least once")
}

func TestWatchdog_NoTampering(t *testing.T) {
	callCount := int32(0)
	var tamperDetected atomic.Bool

	fakeManager := &testutils.FakeManager{
		VerifyLicenseSeatsFn: func(ctx context.Context) error {
			atomic.AddInt32(&callCount, 1)
			return nil // no tampering
		},
	}

	cfg := &server.Config{
		VerifyInterval: 50 * time.Millisecond,
	}

	watchdog := server.NewWatchdog(cfg, fakeManager)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		watchdog.Start(ctx, func() {
			tamperDetected.Store(true)
		})
		close(done)
	}()

	<-done

	assert.False(t, tamperDetected.Load(), "onTamper callback should not have been called")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&callCount), int32(1), "VerifyLicenseSeats should have been called at least once")
}

func TestWatchdog_StopsOnContextCancel(t *testing.T) {
	fakeManager := &testutils.FakeManager{
		VerifyLicenseSeatsFn: func(ctx context.Context) error {
			return nil
		},
	}

	cfg := &server.Config{
		VerifyInterval: 1 * time.Hour, // long interval so it doesn't tick
	}

	watchdog := server.NewWatchdog(cfg, fakeManager)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		watchdog.Start(ctx, func() {})
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// watchdog exited as expected
	case <-time.After(1 * time.Second):
		t.Fatal("watchdog did not stop after context cancellation")
	}
}
