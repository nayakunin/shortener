package config

import "github.com/caarlos0/env/v7"

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var Config config

func LoadConfig() error {
	return env.Parse(&Config)
}
