package config

import (
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/locker"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/server"
)

type Config struct {
	License *licenses.Config
	Server  *server.Config
	Logger  *logger.Config
	Locker  *locker.Config
	DB      *db.Config
}

func New() *Config {
	return &Config{
		Server:  server.NewConfig(),
		License: licenses.NewConfig(),
		Logger:  logger.NewConfig(),
		Locker:  locker.NewConfig(),
		DB:      db.NewConfig(),
	}
}
