package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func Init(config *Config, output io.Writer) {
	var programLevel = new(slog.LevelVar)

	switch config.Verbosity {
	case 0:
		if v := os.Getenv("DEBUG"); v == "true" || v == "t" || v == "1" {
			programLevel.Set(slog.LevelDebug)
		} else {
			programLevel.Set(slog.LevelError)
		}
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

// StringPtr creates an slog attribute for a nullable string pointer.
// If the pointer is nil, it returns a null value, otherwise the dereferenced string.
func StringPtr(key string, val *string) slog.Attr {
	if val == nil {
		return slog.String(key, "<nil>")
	}
	return slog.String(key, *val)
}

// Int64Ptr creates an slog attribute for a nullable int64 pointer.
// If the pointer is nil, it returns a null value, otherwise the dereferenced int64.
func Int64Ptr(key string, val *int64) slog.Attr {
	if val == nil {
		return slog.String(key, "<nil>")
	}
	return slog.Int64(key, *val)
}

// IntPtr creates an slog attribute for a nullable int pointer.
// If the pointer is nil, it returns a null value, otherwise the dereferenced int.
func IntPtr(key string, val *int) slog.Attr {
	if val == nil {
		return slog.String(key, "<nil>")
	}
	return slog.Int(key, *val)
}

// BoolPtr creates an slog attribute for a nullable bool pointer.
// If the pointer is nil, it returns a null value, otherwise the dereferenced bool.
func BoolPtr(key string, val *bool) slog.Attr {
	if val == nil {
		return slog.String(key, "<nil>")
	}
	return slog.Bool(key, *val)
}
