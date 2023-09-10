package grpc

import (
	"context"

	pb "github.com/nayakunin/shortener/proto"
)

func (s *Server) Stats(ctx context.Context, in *pb.Empty) (*pb.StatsReply, error) {
	stats, err := s.Shortener.Stats()
	if err != nil {
		return nil, err
	}

	return &pb.StatsReply{
		Urls:  int64(stats.Urls),
		Users: int64(stats.Users),
	}, nil
}
