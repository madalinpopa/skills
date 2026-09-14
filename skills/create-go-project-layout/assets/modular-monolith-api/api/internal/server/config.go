package server

type Config struct{}

func LoadConfigFromEnv() (Config, error) {
	return Config{}, nil
}
