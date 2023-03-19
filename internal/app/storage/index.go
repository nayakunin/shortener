package storage

import (
	"github.com/pkg/errors"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

var ErrKeyExists = errors.New("key already exists")

type Storager interface {
	Get(key string) (string, bool)
	Add(link string, userID string) (string, error)
	GetUrlsByUser(userID string) (map[string]string, error)
}

func New(cfg config.Config) (Storager, error) {
	if cfg.DatabaseDSN != "" {
		s, err := newDBStorage(cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	if cfg.FileStoragePath != "" {
		s := newFileStorage(cfg.FileStoragePath)
		err := s.restoreData()
		if err != nil {
			return nil, err
		}
		return &s, nil
	}

	s := newStorage()
	return &s, nil
}
