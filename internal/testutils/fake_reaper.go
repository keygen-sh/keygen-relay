package testutils

import (
	"context"

	"github.com/keygen-sh/keygen-relay/internal/server"
)

type FakeReaper struct {
	StartFn   func(ctx context.Context) error
	ManagerFn func() FakeManager
	ConfigFn  func() *server.Config
}

func (r *FakeReaper) Start(ctx context.Context) error {
	if r.StartFn != nil {
		return r.StartFn(ctx)
	}
	return nil
}

func (r *FakeReaper) Manager() FakeManager {
	if r.ManagerFn != nil {
		return r.ManagerFn()
	}

	return FakeManager{}
}

func (r *FakeReaper) Config() *server.Config {
	if r.ConfigFn != nil {
		return r.ConfigFn()
	}

	return &server.Config{}
}
