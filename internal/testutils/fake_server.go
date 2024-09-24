package testutils

import (
	"github.com/gorilla/mux"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/server"
	"sync"
)

type FakeServer struct {
	RunCalled   bool
	RunErr      error
	RunCalledMu sync.Mutex
	ConfigData  *server.Config
	manager     *FakeManager
}

func (s *FakeServer) Run() error {
	s.RunCalledMu.Lock()
	defer s.RunCalledMu.Unlock()
	s.RunCalled = true
	return s.RunErr
}

func (s *FakeServer) Mount(r *mux.Router) {
}

func (s *FakeServer) Config() *server.Config {
	return s.ConfigData
}

func (s *FakeServer) Manager() licenses.Manager {
	return s.manager
}

func NewMockServer(config *server.Config, manager *FakeManager) *FakeServer {
	return &FakeServer{
		ConfigData: config,
		manager:    manager,
	}
}
