// Package storage provides storage for links.
package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nayakunin/shortener/internal/app/interfaces"
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
	AddBatch(batch []interfaces.BatchInput, userID string) ([]interfaces.DBBatchOutput, error)
	GetUrlsByUser(userID string) (map[string]string, error)
	DeleteUserUrls(userID string, keys []string) error
	Stats() (interfaces.Stats, error)
}

// NewStorage returns new storage
func NewStorage(cfg config.Config) (Storager, error) {
	if cfg.DatabaseDSN != "" {
		pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init database")
		}
		rb := newRequestBuffer(MaxRequests)
		s, err := newDBStorage(pool, initDB, rb)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init database")
		}
		return s, nil
	}

	if cfg.FileStoragePath != "" {
		s := newFileStorage(cfg.FileStoragePath)
		err := s.restoreData()
		if err != nil {
			return nil, errors.Wrap(err, "failed to restore data")
		}
		return &s, nil
	}

	s := newStorage()
	return &s, nil
}
