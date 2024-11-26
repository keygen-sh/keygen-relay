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
	if isValidStrategy(v) {
		*e = StrategyType(v)
		return nil
	}
	return errors.New(`must be one of "lifo", "fifo", or "rand"`)
}

func (e *StrategyType) Type() string {
	return "StrategyType"
}

func isValidStrategy(v string) bool {
	switch StrategyType(v) {
	case LIFO, FIFO, RandOrder:
		return true
	default:
		return false
	}
}

type Config struct {
	ServerAddr       string
	ServerPort       int
	EnabledHeartbeat bool
	TTL              time.Duration
	Strategy         StrategyType
	CullInterval     time.Duration
}

func NewConfig() *Config {
	return &Config{
		ServerAddr:       "0.0.0.0",
		ServerPort:       6349,
		TTL:              30 * time.Second,
		EnabledHeartbeat: true,
		Strategy:         FIFO,
		CullInterval:     15 * time.Second,
	}
}
