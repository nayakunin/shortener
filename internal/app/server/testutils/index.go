package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

// MockLink is a mock for storage.Link
type MockLink struct {
	ShortURL    string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

// MockStorage is a mock struct for interfaces.Storage
type MockStorage struct {
	links map[string]MockLink
	users map[string][]MockLink
}

// NewMockStorage creates a new mock storage
func NewMockStorage(initialLinks []MockLink) interfaces.Storage {
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

// Get implements storage.Storager
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

// Add implements storage.Storager
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

// GetUrlsByUser implements interfaces.Storage
func (s *MockStorage) GetUrlsByUser(userID string) (map[string]string, error) {
	links := make(map[string]string)
	for _, link := range s.users[userID] {
		links[link.ShortURL] = link.OriginalURL
	}

	return links, nil
}

// AddBatch implements interfaces.Storage
func (s *MockStorage) AddBatch(batches []interfaces.BatchInput, userID string) ([]interfaces.BatchOutput, error) {
	output := make([]interfaces.BatchOutput, len(batches))
	for i, linkObject := range batches {
		key, err := s.Add(linkObject.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		output[i] = interfaces.BatchOutput{
			Key:           key,
			CorrelationID: linkObject.CorrelationID,
		}
	}
	return output, nil
}

// DeleteUserUrls implements interfaces.Storage
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

// NewMockConfig creates a new mock config
func NewMockConfig() config.Config {
	return config.Config{
		BaseURL: "http://localhost:8080",
	}
}

// AddContext adds config and uuid to gin context
func AddContext(r *gin.Engine, cfg config.Config, userID string) {
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("uuid", userID)
		c.Next()
	})
}
