// Package config provides configuration for the server.
package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

const defaultServerAddress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultFilePath = ""
const defaultDatabaseDSN = ""
const defaultAuthSecret = "secret"

// Config contains configuration for the server.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	AuthSecret      string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	flagsConfig := new(Config)
	flag.StringVar(&flagsConfig.ServerAddress, "a", defaultServerAddress, "server address")
	flag.StringVar(&flagsConfig.BaseURL, "b", defaultBaseURL, "base url")
	flag.StringVar(&flagsConfig.FileStoragePath, "f", defaultFilePath, "file storage path")
	flag.StringVar(&flagsConfig.DatabaseDSN, "d", defaultDatabaseDSN, "database dsn")
	flag.Parse()

	if config.ServerAddress == "" {
		config.ServerAddress = flagsConfig.ServerAddress
	}

	if config.BaseURL == "" {
		config.BaseURL = flagsConfig.BaseURL
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = flagsConfig.FileStoragePath
	}

	if config.DatabaseDSN == "" {
		config.DatabaseDSN = flagsConfig.DatabaseDSN
	}

	if config.AuthSecret == "" {
		config.AuthSecret = defaultAuthSecret
	}

	return &config, nil
}
