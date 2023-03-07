package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

const defaultServerAddress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"

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
	flag.StringVar(&flagsConfig.ServerAddress, "a", "", "server address")
	flag.StringVar(&flagsConfig.BaseURL, "b", "", "base url")
	flag.StringVar(&flagsConfig.FileStoragePath, "f", "", "file storage path")
	flag.Parse()

	if Config.ServerAddress == "" {
		if flagsConfig.ServerAddress == "" {
			flagsConfig.ServerAddress = defaultServerAddress
		}

		Config.ServerAddress = flagsConfig.ServerAddress
	}

	if Config.BaseURL == "" {
		if flagsConfig.BaseURL == "" {
			flagsConfig.BaseURL = defaultBaseURL
		}
		Config.BaseURL = flagsConfig.BaseURL
	}

	if Config.FileStoragePath == "" {
		Config.FileStoragePath = flagsConfig.FileStoragePath
	}

	return nil
}
