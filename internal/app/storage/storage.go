package storage

import (
	"errors"
	"os"
	"sync"

	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/utils"
)

var ErrKeyExists = errors.New("key already exists")

type Storager interface {
	Get(key string) (string, bool)
	Add(link string) (string, error)
}

type Storage struct {
	sync.Mutex
	file  *os.File
	links map[string]string
}

func New() *Storage {
	if config.Config.FileStoragePath == "" {
		return &Storage{
			links: make(map[string]string),
		}
	}

	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var links map[string]string
	if config.Config.FileStoragePath != "" {
		links, err = utils.ReadLinksFromFile(file)
		if err != nil {
			panic(err)
		}
	} else {
		links = make(map[string]string)
	}

	return &Storage{
		file:  file,
		links: links,
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link, ok
}

func (s *Storage) Add(link string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return "", ErrKeyExists
	}

	s.links[key] = link
	if s.file != nil {
		if err := utils.WriteLinkToFile(s.file, key, link); err != nil {
			return "", err
		}
	}

	return key, nil
}

func (s *Storage) Close() error {
	if s.file != nil {
		return s.file.Close()
	}

	return nil
}
