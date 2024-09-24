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

	switch config.Verbosity {
	case 0:
		programLevel.Set(slog.LevelError)
	case 1:
		programLevel.Set(slog.LevelError)
	case 2:
		programLevel.Set(slog.LevelWarn)
	case 3:
		programLevel.Set(slog.LevelInfo)
	default:
		programLevel.Set(slog.LevelDebug)
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
