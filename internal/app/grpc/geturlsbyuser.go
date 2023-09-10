package grpc

import (
	"context"

	"github.com/nayakunin/shortener/internal/app/services/shortener"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUrlsByUser(ctx context.Context, in *pb.GetUrlsByUserRequest) (*pb.GetUrlsByUserReply, error) {
	if in.UserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user id is required")
	}

	urls, err := s.Shortener.GetUrlsByUser(in.UserID)
	if err != nil {
		if errors.Is(err, shortener.ErrNoUrlsFound) {
			return nil, status.Errorf(codes.NotFound, "no urls found")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	result := make([]*pb.GetUrlsByUserLink, len(urls))
	for i, url := range urls {
		result[i] = &pb.GetUrlsByUserLink{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
		}
	}

	return &pb.GetUrlsByUserReply{
		Urls: result,
	}, nil
}
