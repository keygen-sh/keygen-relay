package config

import (
	"github.com/keygen-sh/keygen-relay/internal/db"
	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
	"github.com/keygen-sh/keygen-relay/internal/server"
)

type Config struct {
	License *licenses.Config
	Server  *server.Config
	Logger  *logger.Config
	DB      *db.Config
}

func New() *Config {
	config := Config{
		Server:  server.NewConfig(),
		License: licenses.NewConfig(),
		Logger:  logger.NewConfig(),
		DB:      db.NewConfig(),
	}

	return &config
}
