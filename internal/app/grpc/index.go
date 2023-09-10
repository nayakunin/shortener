package grpc

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
	pb "github.com/nayakunin/shortener/proto"
)

// Shortener is an interface for shortener service
type Shortener interface {
	Shorten(userID string, url string) (string, error)
	ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error)
	Get(key string) (string, error)
	DeleteUserUrls(userID string, keys []string) error
	GetUrlsByUser(userID string) ([]interfaces.Link, error)
	Stats() (*interfaces.Stats, error)
	Ping() error
}

// Server is a struct of the grpc.
type Server struct {
	pb.UnimplementedShortenerServer
	Shortener Shortener
}

// NewServer returns new grpc server
func NewServer(shortener Shortener) *Server {
	return &Server{
		Shortener: shortener,
	}
}
