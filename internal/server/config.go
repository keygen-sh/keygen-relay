package server

type Config struct {
	ServerPort       int
	EnableHeartbeats bool
}

func NewConfig() *Config {
	return &Config{}
}
