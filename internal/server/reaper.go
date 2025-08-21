package server

import (
	"context"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
)

type Reaper interface {
	Start(ctx context.Context) error
	Manager() licenses.Manager
	Config() *Config
}

type reaper struct {
	manager licenses.Manager
	config  *Config
}

func (r *reaper) Start(ctx context.Context) error {
	ticker := time.NewTicker(r.config.CullInterval)
	defer ticker.Stop()

	logger.Debug("starting reaper", "ttl", r.config.TTL, "interval", r.config.CullInterval)

	for {
		select {
		case <-ticker.C:
			r.cull(ctx)
		case <-ctx.Done():
			logger.Debug("stopping reaper")
			return nil
		}
	}
}

func (r *reaper) Manager() licenses.Manager {
	return r.manager
}

func (r *reaper) Config() *Config {
	return r.config
}

func (r *reaper) cull(ctx context.Context) {
	nodes, err := r.manager.CullDeadNodes(ctx, r.config.TTL)
	if err != nil {
		logger.Error("reaper failed to cull dead nodes", "error", err)

		return
	}

	if len(nodes) > 0 {
		logger.Debug("reaper successfully culled dead nodes", "count", len(nodes))
	} else {
		logger.Debug("reaper has nothing to cull")
	}
}

func NewReaper(c *Config, m licenses.Manager) Reaper {
	return &reaper{config: c, manager: m}
}
