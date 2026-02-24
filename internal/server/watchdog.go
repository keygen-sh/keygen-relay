package server

import (
	"context"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/licenses"
	"github.com/keygen-sh/keygen-relay/internal/logger"
)

// Watchdog periodically re-verifies that the seats stored in the database have
// not been tampered with while the server is running. It re-decrypts the
// cryptographically signed license files and compares the maxMachines metadata
// against the seats column. If a discrepancy is found the provided onTamper
// callback is invoked, which typically triggers a graceful server shutdown.
type Watchdog struct {
	manager licenses.Manager
	config  *Config
}

func NewWatchdog(c *Config, m licenses.Manager) *Watchdog {
	return &Watchdog{config: c, manager: m}
}

func (w *Watchdog) Start(ctx context.Context, onTamper func()) {
	ticker := time.NewTicker(w.config.VerifyInterval)
	defer ticker.Stop()

	logger.Debug("starting seat watchdog", "interval", w.config.VerifyInterval)

	for {
		select {
		case <-ticker.C:
			if err := w.manager.VerifyLicenseSeats(ctx); err != nil {
				logger.Error("seat verification failed", "error", err)

				if onTamper != nil {
					onTamper()
				}

				return
			}

			logger.Debug("seat verification passed")
		case <-ctx.Done():
			logger.Debug("stopping seat watchdog")
			return
		}
	}
}
