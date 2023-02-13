package testutils

import "github.com/nayakunin/shortener/internal/app/storage"

type MockStorage struct {
	links map[string]string
}

func NewMockStorage(initialLinks *map[string]string) *MockStorage {
	var links map[string]string
	if initialLinks != nil {
		links = *initialLinks
	} else {
		links = make(map[string]string)
	}

	return &MockStorage{
		links,
	}
}

func (s *MockStorage) Get(key string) (string, bool) {
	link, ok := s.links[key]
	return link, ok
}

func (s *MockStorage) Add(link string) (string, error) {
	key := "link"

	if _, ok := s.links[key]; ok {
		return "", storage.ErrKeyExists
	}

	s.links[key] = link
	return key, nil
}
