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

func (s *Storage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link.LongUrl, ok
}

func (s *Storage) Add(link string, userId string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return "", ErrKeyExists
	}

	linkObject := Link{
		ShortUrl: key,
		LongUrl:  link,
		UserId:   userId,
	}

	s.links[key] = linkObject
	s.users[userId] = append(s.users[userId], linkObject)

	return key, nil
}

func (s *Storage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)
	for _, link := range s.users[id] {
		links[link.ShortUrl] = link.LongUrl
	}

	return links, nil
}
