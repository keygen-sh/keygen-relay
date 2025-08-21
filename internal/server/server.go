package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
)

type Server interface {
	Run() error
	Mount(r *mux.Router)
	Config() *Config
	Manager() licenses.Manager
	Reaper() Reaper
}

type server struct {
	config  *Config
	router  *mux.Router
	manager licenses.Manager
	reaper  Reaper
}

func New(c *Config, m licenses.Manager) Server {
	return &server{
		config:  c,
		router:  mux.NewRouter(),
		manager: m,
		reaper:  NewReaper(c, m),
	}
}

func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := fmt.Sprintf("%s:%d", s.config.ServerAddr, s.config.ServerPort)

	logger.Info("starting server", "addr", s.config.ServerAddr, "port", s.config.ServerPort, "pool", s.config.Pool)

	if s.Config().EnabledHeartbeat {
		go s.reaper.Start(ctx)
	}

	if err := http.ListenAndServe(addr, s.router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed to start", "error", err)

		cancel()

		return err
	}

	logger.Info("server stopped")

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

func (s *server) Reaper() Reaper {
	return s.reaper
}
