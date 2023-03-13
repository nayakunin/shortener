package storage

import (
	"github.com/pkg/errors"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

var ErrKeyExists = errors.New("key already exists")
var ErrFileRead = errors.New("error reading file")

type Storager interface {
	Get(key string) (string, bool)
	Add(link string) (string, error)
}

func newStorage() Storager {
	return &Storage{
		links: make(map[string]string),
	}
}

func newFileStorage(links map[string]string, fileStoragePath string) Storager {
	return &FileStorage{
		Storage: Storage{
			links: links,
		},
		fileStoragePath: fileStoragePath,
	}
}

func New(cfg config.Config) (Storager, error) {
	if cfg.FileStoragePath == "" {
		return newStorage(), nil
	}

	links, err := restoreLinksFromFile(cfg.FileStoragePath)
	if err != nil {
		return nil, ErrFileRead
	}

	return newFileStorage(links, cfg.FileStoragePath), nil
}
