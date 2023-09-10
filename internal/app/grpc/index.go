package grpc

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
	pb "github.com/nayakunin/shortener/proto"
)

type Shortener interface {
	Shorten(userID string, url string) (string, error)
	ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error)
	Get(key string) (string, error)
	DeleteUserUrls(userID string, keys []string) error
	GetUrlsByUser(userID string) ([]interfaces.Link, error)
	Stats() (*interfaces.Stats, error)
	Ping() error
}

type Server struct {
	pb.UnimplementedShortenerServer
	Shortener Shortener
}

func NewServer(shortener Shortener) *Server {
	return &Server{
		Shortener: shortener,
	}
}
