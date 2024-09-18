package logger

type Config struct {
	Verbosity int
}

func NewConfig() *Config {
	return &Config{Verbosity: 0}
}
