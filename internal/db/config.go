package db

type Config struct {
	DatabaseFilePath string
}

func NewConfig() *Config {
	return &Config{}
}
