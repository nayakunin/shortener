package testutils

import "github.com/nayakunin/shortener/internal/app/config"

// NewMockConfig creates a new mock config
func NewMockConfig() config.Config {
	return config.Config{
		BaseURL: "http://localhost:8080",
	}
}
