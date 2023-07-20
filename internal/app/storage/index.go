// Package storage provides storage for links.
package storage

import (
	"github.com/pkg/errors"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

// ErrKeyExists is returned when key already exists
var ErrKeyExists = errors.New("key already exists")

// ErrBatchInvalidURL is returned when url is invalid
var ErrBatchInvalidURL = errors.New("invalid url")

// ErrKeyDeleted is returned when key is deleted
var ErrKeyDeleted = errors.New("key deleted")

// ErrKeyNotFound is returned when key is not found
var ErrKeyNotFound = errors.New("key not found")

// Storager is an interface for storage
type Storager interface {
	Get(key string) (string, error)
	Add(link string, userID string) (string, error)
	AddBatch(batch []BatchInput, userID string) ([]BatchOutput, error)
	GetUrlsByUser(userID string) (map[string]string, error)
	DeleteUserUrls(userID string, keys []string) error
}

// New returns new storage
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
