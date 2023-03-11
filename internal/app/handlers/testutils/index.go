package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

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
