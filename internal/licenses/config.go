package licenses

import (
	"time"
)

type Config struct {
	Strategy          string
	TTL               time.Duration
	EnabledAudit      bool
	ExtendOnHeartbeat bool
}

func NewConfig() *Config {
	return &Config{Strategy: "fifo"}
}
