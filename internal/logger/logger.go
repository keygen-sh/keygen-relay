package logger

import (
	"io"
	"log/slog"
	"os"
	"reflect"
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

// Debug logs at debug level with automatic nullable type handling
func Debug(msg string, args ...any) {
	slog.Debug(msg, processArgs(args...)...)
}

// Info logs at info level with automatic nullable type handling
func Info(msg string, args ...any) {
	slog.Info(msg, processArgs(args...)...)
}

// Warn logs at warn level with automatic nullable type handling
func Warn(msg string, args ...any) {
	slog.Warn(msg, processArgs(args...)...)
}

// Error logs at error level with automatic nullable type handling
func Error(msg string, args ...any) {
	slog.Error(msg, processArgs(args...)...)
}

// derefPointer returns "<nil>" for nil pointers, the pointed-to value for non-nil
// pointers, and the original value for non-pointers.
func derefPointer(v any) any {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() { // nil interface
		return v
	}

	if rv.Kind() != reflect.Ptr {
		return v
	}

	if rv.IsNil() {
		return "<nil>"
	}

	return rv.Elem().Interface()
}

// processArgs processes alternating key/value args, dereferencing *any pointer type.
func processArgs(in ...any) []any {
	out := make([]any, 0, len(in))

	for i := 0; i < len(in); i += 2 {
		out = append(out, in[i]) // key

		if i+1 < len(in) { // value
			out = append(out, derefPointer(in[i+1]))
		}
	}

	return out
}
