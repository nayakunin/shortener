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

func (s *MockStorage) Get(key string) (string, bool) {
	link, ok := s.links[key]
	return link.OriginalURL, ok
}

func (s *MockStorage) Add(link string, userID string) (string, error) {
	key := "link"

	if _, ok := s.links[key]; ok {
		return "", storage.ErrKeyExists
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

func NewMockConfig() config.Config {
	return config.Config{
		BaseURL: "http://localhost:8080",
	}
}

func AddContext(r *gin.Engine, cfg config.Config) {
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})
}
