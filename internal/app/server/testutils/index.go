package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type MockLink struct {
	ShortUrl string
	LongUrl  string
	UserId   string
}

type MockStorage struct {
	links map[string]MockLink
	users map[string][]MockLink
}

func NewMockStorage(initialLinks []MockLink) *MockStorage {
	var links map[string]MockLink
	var users map[string][]MockLink
	if initialLinks != nil {
		links = make(map[string]MockLink)
		for _, link := range initialLinks {
			links[link.ShortUrl] = link
			users[link.UserId] = append(users[link.UserId], link)
		}
	} else {
		links = make(map[string]MockLink)
		users = make(map[string][]MockLink)
	}

	return &MockStorage{
		links: links,
		users: users,
	}
}

func (s *MockStorage) Get(key string) (string, bool) {
	link, ok := s.links[key]
	return link.LongUrl, ok
}

func (s *MockStorage) Add(link string, userId string) (string, error) {
	key := "link"

	if _, ok := s.links[key]; ok {
		return "", storage.ErrKeyExists
	}

	linkObject := MockLink{
		ShortUrl: key,
		LongUrl:  link,
		UserId:   userId,
	}

	s.links[key] = linkObject
	s.users[userId] = append(s.users[userId], linkObject)
	return key, nil
}

func (s *MockStorage) GetUrlsByUser(userId string) (map[string]string, error) {
	links := make(map[string]string)
	for _, link := range s.users[userId] {
		links[link.ShortUrl] = link.LongUrl
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
