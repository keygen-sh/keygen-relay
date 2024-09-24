package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/keygen-sh/keygen-relay/internal/cmd"
	"github.com/keygen-sh/keygen-relay/internal/config"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"github.com/keygen-sh/keygen-relay/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServeCmd_Defaults(t *testing.T) {
	cfg := &server.Config{
		ServerPort:       8080,
		TTL:              30 * time.Second,
		EnabledHeartbeat: true,
		Strategy:         server.FIFO,
	}

	mockServer := testutils.NewMockServer(cfg, &testutils.FakeManager{})

	serveCmd := cmd.ServeCmd(mockServer)

	serveCmd.SetArgs([]string{})

	output := &bytes.Buffer{}
	serveCmd.SetOut(output)

	err := serveCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "The server is starting")
	assert.True(t, mockServer.RunCalled)
	assert.Equal(t, 8080, cfg.ServerPort)
	assert.Equal(t, 30*time.Second, cfg.TTL)
	assert.True(t, cfg.EnabledHeartbeat)
	assert.Equal(t, server.FIFO, cfg.Strategy)
}

func TestServeCmd_Flags(t *testing.T) {
	cfg := &config.Config{
		Server: &server.Config{
			ServerPort:       8080,
			TTL:              30 * time.Second,
			EnabledHeartbeat: true,
			Strategy:         server.FIFO,
			CleanupInterval:  5 * time.Second,
		},
		License: &licenses.Config{},
	}

	manager := testutils.FakeManager{
		ConfigFn: func() *licenses.Config {
			return cfg.License
		},
	}

	mockServer := testutils.NewMockServer(cfg.Server, &manager)

	serveCmd := cmd.ServeCmd(mockServer)

	serveCmd.SetArgs([]string{
		"--port", "9090",
		"--ttl", "1m",
		"--no-heartbeats",
		"--strategy", "lifo",
	})

	output := &bytes.Buffer{}
	serveCmd.SetOut(output)

	err := serveCmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "The server is starting")
	assert.True(t, mockServer.RunCalled)
	assert.Equal(t, 9090, cfg.Server.ServerPort)
	assert.Equal(t, 1*time.Minute, cfg.Server.TTL)
	assert.False(t, cfg.Server.EnabledHeartbeat)
	assert.Equal(t, server.LIFO, cfg.Server.Strategy)
	assert.Equal(t, string(cfg.Server.Strategy), cfg.License.Strategy)
	assert.Equal(t, cfg.Server.EnabledHeartbeat, cfg.License.ExtendOnHeartbeat)
	assert.Equal(t, cfg.Server.CleanupInterval, 5*time.Second)
}

func TestServeCmd_InvalidStrategy(t *testing.T) {
	cfg := &config.Config{
		Server: &server.Config{
			ServerPort:       8080,
			TTL:              30 * time.Second,
			EnabledHeartbeat: true,
			Strategy:         server.FIFO,
		},
	}

	mockServer := testutils.NewMockServer(cfg.Server, &testutils.FakeManager{})
	serveCmd := cmd.ServeCmd(mockServer)

	serveCmd.SetArgs([]string{
		"--strategy", "invalid",
	})

	output := &bytes.Buffer{}
	serveCmd.SetOut(output)
	serveCmd.SetErr(output)

	err := serveCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid argument \"invalid\" for \"--strategy\"")
	assert.False(t, mockServer.RunCalled)
}

func TestServeCmd_RunError(t *testing.T) {
	cfg := &config.Config{
		Server: &server.Config{
			ServerPort:       8080,
			TTL:              30 * time.Second,
			EnabledHeartbeat: true,
			Strategy:         server.FIFO,
		},
	}

	mockServer := testutils.NewMockServer(cfg.Server, &testutils.FakeManager{})
	mockServer.RunErr = fmt.Errorf("failed to start server")

	serveCmd := cmd.ServeCmd(mockServer)

	output := &bytes.Buffer{}
	serveCmd.SetOut(output)
	serveCmd.SetErr(output)

	err := serveCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error running server")
	assert.True(t, mockServer.RunCalled)
}
