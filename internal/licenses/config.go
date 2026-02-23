package licenses

import (
	"time"
)

type Config struct {
	Strategy          string
	TTL               time.Duration
	EnabledAudit      bool
	ExtendOnHeartbeat bool
	PublicKey         string
}

func NewConfig() *Config {
	return &Config{Strategy: "fifo"}
}
