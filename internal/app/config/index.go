// Package config provides configuration for the rest.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/pkg/errors"
)

const defaultServerAddress = "localhost:8080"
const defaultBaseURL = "http://localhost:8080"
const defaultFilePath = ""
const defaultDatabaseDSN = ""
const defaultAuthSecret = "secret"

// Config contains configuration for the rest.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL         string `env:"BASE_URL" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	Config          string `env:"CONFIG"`
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
	flag.StringVar(&flagsConfig.ServerAddress, "a", defaultServerAddress, "rest address")
	flag.StringVar(&flagsConfig.BaseURL, "b", defaultBaseURL, "base url")
	flag.StringVar(&flagsConfig.FileStoragePath, "f", defaultFilePath, "file storage path")
	flag.StringVar(&flagsConfig.DatabaseDSN, "d", defaultDatabaseDSN, "database dsn")
	flag.BoolVar(&flagsConfig.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&flagsConfig.TrustedSubnet, "t", "", "trusted subnet")
	flag.StringVar(&flagsConfig.Config, "c", "", "config file")
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

	if !config.EnableHTTPS {
		config.EnableHTTPS = flagsConfig.EnableHTTPS
	}

	if config.AuthSecret == "" {
		config.AuthSecret = defaultAuthSecret
	}

	if config.Config == "" {
		config.Config = flagsConfig.Config
	}

	// read config from a JSON file
	if config.Config != "" {
		fileBytes, err := os.ReadFile(config.Config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config file")
		}

		if err := json.Unmarshal(fileBytes, &config); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal config file")
		}
	}

	return &config, nil
}
