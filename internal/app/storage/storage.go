package storage

import (
	"github.com/nayakunin/shortener/internal/app/utils"
	"sync"
)

type Storage struct {
	mu    sync.Mutex
	links map[string]string
}

func New() *Storage {
	return &Storage{
		links: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, ok := s.links[key]
	return link, ok
}

func (s *Storage) Add(link string) string {
	key := utils.RandSeq(5)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.links[key] = link

	return key
}
