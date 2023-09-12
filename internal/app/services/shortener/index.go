package shortener

import (
	"github.com/nayakunin/shortener/internal/app/config"
	"github.com/nayakunin/shortener/internal/app/interfaces"
)

// Shortener is an interface for shortener service
type Shortener interface {
	Shorten(userID string, url string) (string, error)
	ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error)
	Get(key string) (string, error)
	DeleteUserUrls(userID string, keys []string) error
	GetUrlsByUser(userID string) ([]interfaces.Link, error)
	Stats() (*interfaces.Stats, error)
	Ping() error
}

// Service is a struct of the shortener.
type Service struct {
	Cfg     config.Config
	Storage interfaces.Storage
}

// NewShortenerService is a constructor for the shortener service.
func NewShortenerService(cfg config.Config, s interfaces.Storage) *Service {
	return &Service{
		Cfg:     cfg,
		Storage: s,
	}
}
