package storage

import (
	"github.com/nayakunin/shortener/internal/app/utils"
)

type FileStorage struct {
	Storage
	fileStoragePath string
}

func (s *FileStorage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link, ok
}

func (s *FileStorage) Add(link string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return "", ErrKeyExists
	}

	s.links[key] = link
	if err := writeLinkToFile(s.fileStoragePath, key, link); err != nil {
		return "", err
	}

	return key, nil
}
