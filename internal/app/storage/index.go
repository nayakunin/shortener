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

func newStorage() Storage {
	return Storage{
		links: make(map[string]Link),
		users: make(map[string][]Link),
	}
}

func newFileStorage(fileStoragePath string) FileStorage {
	return FileStorage{
		Storage: Storage{
			links: make(map[string]Link),
			users: make(map[string][]Link),
		},
		fileStoragePath: fileStoragePath,
	}
}

func New(cfg config.Config) (Storager, error) {
	if cfg.FileStoragePath == "" {
		s := newStorage()
		return &s, nil
	}

	s := newFileStorage(cfg.FileStoragePath)
	err := s.restoreData()
	if err != nil {
		return nil, err
	}

	return &s, nil
}
