package logger

import (
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"io"
	"log/slog"
	"os"
	"time"
)

func Init(config *Config, output io.Writer) {
	var programLevel = new(slog.LevelVar)

	switch {
	case config.Verbosity == 0:
		programLevel.Set(slog.LevelError)
	case config.Verbosity == 1:
		programLevel.Set(slog.LevelError)
	case config.Verbosity == 2:
		programLevel.Set(slog.LevelWarn)
	case config.Verbosity == 3:
		programLevel.Set(slog.LevelInfo)
	case config.Verbosity >= 4:
		programLevel.Set(slog.LevelDebug)
	default:
		programLevel.Set(slog.LevelError)
	}

	handler := tint.NewHandler(output, &tint.Options{
		Level:      programLevel,
		TimeFormat: time.DateTime,
		NoColor:    config.DisableColor || !isatty.IsTerminal(os.Stdout.Fd()),
		AddSource:  true,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
