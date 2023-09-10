package grpc

import (
	"context"
	"errors"

	"github.com/nayakunin/shortener/internal/app/storage"
	pb "github.com/nayakunin/shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetLink returns link by key
func (s *Server) GetLink(ctx context.Context, in *pb.GetLinkRequest) (*pb.GetLinkReply, error) {
	if in.Key == "" {
		return nil, status.Errorf(codes.InvalidArgument, "key is required")
	}

	link, err := s.Shortener.Get(in.Key)
	if err != nil {
		if errors.Is(err, storage.ErrKeyDeleted) {
			return nil, status.Errorf(codes.NotFound, "key is deleted")
		}

		return nil, status.Errorf(codes.NotFound, "key not found")
	}

	return &pb.GetLinkReply{
		Url: link,
	}, nil
}
