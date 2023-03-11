package storage

import (
	"errors"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

var ErrKeyExists = errors.New("key already exists")

type Storager interface {
	Get(key string) (string, bool)
	Add(link string) (string, error)
}

func New(cfg config.Config) (Storager, error) {
	if cfg.FileStoragePath == "" {
		return &Storage{
			links: make(map[string]string),
		}, nil
	}

	return restoreLinksFromFile(cfg.FileStoragePath)
}
