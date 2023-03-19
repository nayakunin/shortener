package storage

import (
	"sync"

	"github.com/nayakunin/shortener/internal/app/utils"
)

type Storage struct {
	sync.Mutex
	links map[string]Link
	users map[string][]Link
}

func newStorage() Storage {
	return Storage{
		links: make(map[string]Link),
		users: make(map[string][]Link),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link.OriginalURL, ok
}

func (s *Storage) Add(link string, userID string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return "", ErrKeyExists
	}

	linkObject := Link{
		ShortURL:    key,
		OriginalURL: link,
		UserID:      userID,
	}

	s.links[key] = linkObject
	s.users[userID] = append(s.users[userID], linkObject)

	return key, nil
}

func (s *Storage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)
	for _, link := range s.users[id] {
		links[link.ShortURL] = link.OriginalURL
	}

	return links, nil
}
