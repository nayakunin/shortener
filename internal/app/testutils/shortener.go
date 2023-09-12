package testutils

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
)

// MockShortenerServiceParameters is a mock struct for interfaces.ShortenerService
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

// MockShortenerService is a mock struct for interfaces.ShortenerService
type MockShortenerService struct {
	parameters MockShortenerServiceParameters
}

// NewMockShortenerService creates a new mock storage
func NewMockShortenerService(parameters MockShortenerServiceParameters) *MockShortenerService {
	return &MockShortenerService{
		parameters: parameters,
	}
}

// Shorten implements interfaces.ShortenerService
func (s *MockShortenerService) Shorten(userID string, url string) (string, error) {
	return s.parameters.ShortenReplyString, s.parameters.ShortenReplyError
}

// ShortenBatch implements interfaces.ShortenerService
func (s *MockShortenerService) ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error) {
	return s.parameters.ShortenBatchReply, s.parameters.ShortenBatchError
}

// Get implements interfaces.ShortenerService
func (s *MockShortenerService) Get(key string) (string, error) {
	return s.parameters.GetReplyString, s.parameters.GetReplyError
}

// DeleteUserUrls implements interfaces.ShortenerService
func (s *MockShortenerService) DeleteUserUrls(userID string, keys []string) error {
	return s.parameters.DeleteReplyError
}

// GetUrlsByUser implements interfaces.ShortenerService
func (s *MockShortenerService) GetUrlsByUser(userID string) ([]interfaces.Link, error) {
	return s.parameters.GetUrlsReply, s.parameters.GetUrlsError
}

// Stats implements interfaces.ShortenerService
func (s *MockShortenerService) Stats() (*interfaces.Stats, error) {
	return s.parameters.StatsReply, s.parameters.StatsError
}

// Ping implements interfaces.ShortenerService
func (s *MockShortenerService) Ping() error {
	return s.parameters.PingError
}
