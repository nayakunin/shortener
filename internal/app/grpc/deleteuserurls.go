package grpc

import (
	"context"

	pb "github.com/nayakunin/shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUserUrls deletes urls by keys for user
func (s *Server) DeleteUserUrls(ctx context.Context, in *pb.DeleteUserUrlsRequest) (*pb.Empty, error) {
	if in.UserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user id is required")
	}

	if err := s.Shortener.DeleteUserUrls(in.UserID, in.Keys); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}
