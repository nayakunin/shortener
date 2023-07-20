package storage

import (
	"sync"

	"github.com/nayakunin/shortener/internal/app/utils"
	"github.com/pkg/errors"
)

// Storage is a storage
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

// Get returns original url by key
func (s *Storage) Get(key string) (string, error) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]

	if !ok {
		return "", ErrKeyNotFound
	}

	if link.IsDeleted {
		return "", ErrKeyDeleted
	}

	return link.OriginalURL, nil
}

// Add adds new link to storage
func (s *Storage) Add(link string, userID string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return key, ErrKeyExists
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

// AddBatch adds new links to storage
func (s *Storage) AddBatch(batches []BatchInput, userID string) ([]BatchOutput, error) {
	output := make([]BatchOutput, len(batches))
	for i, linkObject := range batches {
		key, err := s.Add(linkObject.OriginalURL, userID)
		if err != nil && !errors.Is(err, ErrKeyExists) {
			return nil, err
		}
		output[i] = BatchOutput{
			Key:           key,
			CorrelationID: linkObject.CorrelationID,
		}
	}

	return output, nil
}

// GetUrlsByUser returns all urls by user id
func (s *Storage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)
	for _, link := range s.users[id] {
		links[link.ShortURL] = link.OriginalURL
	}

	return links, nil
}

// DeleteUserUrls deletes user's urls
func (s *Storage) DeleteUserUrls(userID string, keys []string) error {
	s.Lock()
	defer s.Unlock()

	userLinks := s.users[userID]

	for _, key := range keys {
		link := s.links[key]
		if link.UserID != userID {
			continue
		}

		s.links[key] = Link{
			ShortURL:    key,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   true,
		}

		for i, userLink := range userLinks {
			if userLink.ShortURL == key {
				userLinks[i] = Link{
					ShortURL:    key,
					OriginalURL: userLink.OriginalURL,
					UserID:      userLink.UserID,
					IsDeleted:   true,
				}
			}
		}
	}

	s.users[userID] = userLinks

	return nil
}
