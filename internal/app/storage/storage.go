package storage

import (
	"sync"

	"github.com/nayakunin/shortener/internal/app/utils"
)

type Storage struct {
	sync.Mutex
	links map[string]string
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
