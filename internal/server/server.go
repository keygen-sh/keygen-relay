package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
)

type Server interface {
	Run() error
	Mount(r *mux.Router)
	Config() *Config
	Manager() licenses.Manager
}

type server struct {
	config  *Config
	router  *mux.Router
	manager licenses.Manager
}

func New(c *Config, m licenses.Manager) Server {
	return &server{
		config:  c,
		router:  mux.NewRouter(),
		manager: m,
	}
}

func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := fmt.Sprintf(":%d", s.config.ServerPort)
	slog.Info("Starting server", "port", s.config.ServerPort)

	if s.Config().EnabledHeartbeat {
		go s.startCullRoutine(ctx)
	}

	err := http.ListenAndServe(address, s.router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed to start", "error", err)
		cancel()
		return err
	}

	slog.Info("Server stopped")
	return nil
}

func (s *server) Mount(r *mux.Router) {
	s.router = r
}

func (s *server) Router() *mux.Router {
	return s.router
}

func (s *server) Config() *Config {
	return s.config
}

func (s *server) Manager() licenses.Manager {
	return s.manager
}

func (s *server) startCullRoutine(ctx context.Context) {
	ticker := time.NewTicker(s.config.CullInterval)
	defer ticker.Stop()

	slog.Debug("Starting cull routine for inactive nodes", "ttl", s.config.TTL, "cullInterval", s.config.CullInterval)

	for {
		select {
		case <-ticker.C:
			s.cull()
		case <-ctx.Done():
			slog.Debug("Stopping cull routine")
			return
		}
	}
}

func (s *server) cull() {
	ctx := context.Background()
	err := s.manager.CullInactiveNodes(ctx, s.Config().TTL)
	if err != nil {
		slog.Error("Failed to cull inactive nodes", "error", err)
	} else {
		slog.Debug("Successfully culled inactive nodes")
	}
}
