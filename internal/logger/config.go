package logger

type Config struct {
	Verbosity    int
	DisableColor bool
}

func NewConfig() *Config {
	return &Config{}
}
