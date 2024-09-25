package logger

import (
	"io"
	"log/slog"
)

func Init(config *Config, output io.Writer) {
	var programLevel = new(slog.LevelVar)

	switch config.Verbosity {
	case 0:
		programLevel.Set(slog.LevelInfo)
	case 1:
		programLevel.Set(slog.LevelWarn)
	case 2:
		programLevel.Set(slog.LevelDebug)
	case 3:
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelError)
	}

	logger := slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: programLevel,
	}))

	slog.SetDefault(logger)
}
