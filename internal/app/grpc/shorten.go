package grpc

import (
	"context"

	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/storage"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Shorten shortens url for user
func (s *Server) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenReply, error) {
	if in.UserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user_id is required")
	}

	shortURL, err := s.Shortener.Shorten(in.UserID, in.Url)
	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid url")
		}

		if errors.Is(err, storage.ErrKeyExists) {
			return &pb.ShortenReply{
				Url: shortURL,
			}, status.Error(codes.AlreadyExists, "url already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.ShortenReply{
		Url: shortURL,
	}, nil
}
