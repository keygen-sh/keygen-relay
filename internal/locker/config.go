package locker

type Config struct {
	MachineFilePath string
	LicenseKey      string
}

func NewConfig() *Config {
	return &Config{}
}
