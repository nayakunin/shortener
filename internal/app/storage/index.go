package storage

import (
	"github.com/pkg/errors"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

var ErrKeyExists = errors.New("key already exists")
var ErrBatchInvalidURL = errors.New("invalid url")
var ErrKeyDeleted = errors.New("key deleted")
var ErrKeyNotFound = errors.New("key not found")

type Storager interface {
	Get(key string) (string, error)
	Add(link string, userID string) (string, error)
	AddBatch(batch []BatchInput, userID string) ([]BatchOutput, error)
	GetUrlsByUser(userID string) (map[string]string, error)
	DeleteUserUrls(userID string, keys []string) error
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
