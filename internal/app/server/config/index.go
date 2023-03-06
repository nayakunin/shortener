package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
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
		if flagsConfig.ServerAddress == "" {
			flagsConfig.ServerAddress = "localhost:8080"
		}

		Config.ServerAddress = flagsConfig.ServerAddress
	}

	if Config.BaseURL == "" {
		if flagsConfig.BaseURL == "" {
			flagsConfig.BaseURL = "http://localhost:8080"
		}
		Config.BaseURL = flagsConfig.BaseURL
	}

	if Config.FileStoragePath == "" {
		Config.FileStoragePath = flagsConfig.FileStoragePath
	}

	return nil
}
