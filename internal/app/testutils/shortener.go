package testutils

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
)

type MockShortenerServiceParameters struct {
	ShortenReplyString string
	ShortenReplyError  error
	ShortenBatchReply  []interfaces.BatchOutput
	ShortenBatchError  error
	GetReplyString     string
	GetReplyError      error
	DeleteReplyError   error
	GetUrlsReply       []interfaces.Link
	GetUrlsError       error
	StatsReply         *interfaces.Stats
	StatsError         error
	PingError          error
}

type MockShortenerService struct {
	parameters MockShortenerServiceParameters
}

func NewMockShortenerService(parameters MockShortenerServiceParameters) shortener.Shortener {
	return &MockShortenerService{
		parameters: parameters,
	}
}

func (s *MockShortenerService) Shorten(userID string, url string) (string, error) {
	return s.parameters.ShortenReplyString, s.parameters.ShortenReplyError
}

func (s *MockShortenerService) ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error) {
	return s.parameters.ShortenBatchReply, s.parameters.ShortenBatchError
}

func (s *MockShortenerService) Get(key string) (string, error) {
	return s.parameters.GetReplyString, s.parameters.GetReplyError
}

func (s *MockShortenerService) DeleteUserUrls(userID string, keys []string) error {
	return s.parameters.DeleteReplyError
}

func (s *MockShortenerService) GetUrlsByUser(userID string) ([]interfaces.Link, error) {
	return s.parameters.GetUrlsReply, s.parameters.GetUrlsError
}

func (s *MockShortenerService) Stats() (*interfaces.Stats, error) {
	return s.parameters.StatsReply, s.parameters.StatsError
}

func (s *MockShortenerService) Ping() error {
	return s.parameters.PingError
}
