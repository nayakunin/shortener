package shortener

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/server/config"
)

type Shortener interface {
	Shorten(userID string, url string) (string, error)
	ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error)
	Get(key string) (string, error)
	DeleteUserUrls(userID string, keys []string) error
	GetUrlsByUser(userID string) ([]interfaces.Link, error)
	Stats() (*interfaces.Stats, error)
	Ping() error
}

type Service struct {
	Cfg     config.Config
	Storage interfaces.Storage
}

func NewShortenerService(cfg config.Config, s interfaces.Storage) *Service {
	return &Service{
		Cfg:     cfg,
		Storage: s,
	}
}
