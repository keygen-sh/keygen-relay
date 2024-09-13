package licenses

import (
	"time"
)

type StrategyType string

const (
	LIFO      StrategyType = "lifo"
	FIFO      StrategyType = "fifo"
	RandOrder StrategyType = "rand"
)

type Config struct {
	Strategy     StrategyType
	TTL          time.Duration
	EnabledAudit bool
}

func NewConfig() *Config {
	return &Config{}
}
