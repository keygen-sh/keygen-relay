package server

import (
	"errors"
	"time"
)

type StrategyType string

const (
	LIFO      StrategyType = "lifo"
	FIFO      StrategyType = "fifo"
	RandOrder StrategyType = "rand"
)

func (e *StrategyType) String() string {
	return string(*e)
}

func (e *StrategyType) Set(v string) error {
	switch v {
	case "lifo", "fifo", "rand":
		*e = StrategyType(v)
		return nil
	default:
		return errors.New(`must be one of "lifo", "fifo", or "rand"`)
	}
}

func (e *StrategyType) Type() string {
	return "StrategyType"
}

type Config struct {
	ServerPort       int
	EnabledHeartbeat bool
	TTL              time.Duration
	Strategy         StrategyType
	CleanupInterval  time.Duration
}

func NewConfig() *Config {
	return &Config{
		ServerPort:       8080,
		TTL:              30 * time.Second,
		EnabledHeartbeat: true,
		Strategy:         FIFO,
		CleanupInterval:  15,
	}
}
