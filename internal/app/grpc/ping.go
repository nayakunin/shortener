package grpc

import (
	"context"

	pb "github.com/nayakunin/shortener/proto"
)

func (s *Server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	err := s.Shortener.Ping()
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
