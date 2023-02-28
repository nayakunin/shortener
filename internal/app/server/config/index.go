package config

import "github.com/caarlos0/env/v7"

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost"`
	Port          string `env:"PORT" envDefault:"8080"`
}

var Config config

func LoadConfig() error {
	return env.Parse(&Config)
}
