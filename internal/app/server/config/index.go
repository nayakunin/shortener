package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var Config config

func LoadConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	flagsConfig := new(config)
	flag.StringVar(&flagsConfig.ServerAddress, "a", Config.ServerAddress, "server address")
	flag.StringVar(&flagsConfig.BaseURL, "b", Config.BaseURL, "base url")
	flag.StringVar(&flagsConfig.FileStoragePath, "f", Config.FileStoragePath, "file storage path")
	flag.Parse()

	if Config.ServerAddress == "" {
		Config.ServerAddress = flagsConfig.ServerAddress
	}

	if Config.BaseURL == "" {
		Config.BaseURL = flagsConfig.BaseURL
	}

	if Config.FileStoragePath == "" {
		Config.FileStoragePath = flagsConfig.FileStoragePath
	}

	return nil
}
