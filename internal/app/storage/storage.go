package storage

import (
	"errors"
	"sync"

	"github.com/nayakunin/shortener/internal/app/utils"
)

var ErrKeyExists = errors.New("key already exists")

type Storager interface {
	Get(key string) (string, bool)
	Add(link string) (string, error)
}

type Storage struct {
	sync.Mutex
	links map[string]string
}

func New() *Storage {
	return &Storage{
		links: make(map[string]string),
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

	return key, nil
}
