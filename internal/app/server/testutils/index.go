package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type MockLink struct {
	ShortURL    string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

type MockStorage struct {
	links map[string]MockLink
	users map[string][]MockLink
}

func NewMockStorage(initialLinks []MockLink) *MockStorage {
	links := make(map[string]MockLink)
	users := make(map[string][]MockLink)

	if len(initialLinks) > 0 {
		for _, link := range initialLinks {
			links[link.ShortURL] = link
			users[link.UserID] = append(users[link.UserID], link)
		}
	}

	return &MockStorage{
		links: links,
		users: users,
	}
}

func (s *MockStorage) Get(key string) (string, error) {
	link, ok := s.links[key]
	if !ok {
		return "", storage.ErrKeyNotFound
	}

	if link.IsDeleted {
		return "", storage.ErrKeyDeleted
	}

	return link.OriginalURL, nil
}

func (s *MockStorage) Add(link string, userID string) (string, error) {
	key := "link"

	if _, ok := s.links[key]; ok {
		return key, storage.ErrKeyExists
	}

	linkObject := MockLink{
		ShortURL:    key,
		OriginalURL: link,
		UserID:      userID,
	}

	s.links[key] = linkObject
	s.users[userID] = append(s.users[userID], linkObject)
	return key, nil
}

func (s *MockStorage) GetUrlsByUser(userID string) (map[string]string, error) {
	links := make(map[string]string)
	for _, link := range s.users[userID] {
		links[link.ShortURL] = link.OriginalURL
	}

	return links, nil
}

func (s *MockStorage) AddBatch(batches []storage.BatchInput, userID string) ([]storage.BatchOutput, error) {
	output := make([]storage.BatchOutput, len(batches))
	for i, linkObject := range batches {
		key, err := s.Add(linkObject.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		output[i] = storage.BatchOutput{
			Key:           key,
			CorrelationID: linkObject.CorrelationID,
		}
	}
	return output, nil
}

func (s *MockStorage) DeleteUserUrls(userID string, keys []string) error {
	userLinks := s.users[userID]

	for _, key := range keys {
		link := s.links[key]
		if link.UserID != userID {
			continue
		}

		delete(s.links, key)
		s.links[key] = MockLink{
			ShortURL:    key,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   true,
		}

		for i, userLink := range userLinks {
			if userLink.ShortURL == key {
				userLinks[i] = MockLink{
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

func NewMockConfig() config.Config {
	return config.Config{
		BaseURL: "http://localhost:8080",
	}
}

func AddContext(r *gin.Engine, cfg config.Config, userID string) {
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("uuid", userID)
		c.Next()
	})
}
