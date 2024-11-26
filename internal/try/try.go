package try

import (
	"os"
	"strconv"
	"time"
)

type Accessor[T comparable] func() T

func Try[T comparable](accessors ...Accessor[T]) T {
	var zero T

	for _, accessor := range accessors {
		if v := accessor(); v != zero {
			return v
		}
	}

	return zero
}

// TODO(ezekg) implement relay toml/yaml config file?
func Config[T comparable](key string) func() T {
	return func() T {
		var zero T

		return zero
	}
}

func Env(key string) func() string {
	return func() string {
		return os.Getenv(key)
	}
}

func EnvAs[T any](key string, converter func(string) T) func() T {
	return func() T {
		value := os.Getenv(key)

		return converter(value)
	}
}

func EnvBool(key string) func() bool {
	return EnvAs(key, func(value string) bool {
		return value == "true" || value == "t" || value == "1"
	})
}

func EnvInt(key string) func() int {
	return EnvAs(key, func(value string) int {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		} else {
			return 0
		}
	})
}

func EnvDuration(key string) func() time.Duration {
	return EnvAs(key, func(value string) time.Duration {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		} else {
			return 0
		}
	})
}

func Static[T comparable](value T) Accessor[T] {
	return func() T {
		return value
	}
}
