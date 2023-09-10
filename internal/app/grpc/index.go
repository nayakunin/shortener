package grpc

import (
	"context"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/server/config"
	pb "github.com/nayakunin/shortener/proto"
)

type Server struct {
	pb.UnimplementedShortenerServer
	Cfg     config.Config
	Storage interfaces.Storage
}

func NewServer(cfg config.Config, s interfaces.Storage) *Server {
	return &Server{
		Cfg:     cfg,
		Storage: s,
	}
}

func (s *Server) DeleteUserUrls(ctx context.Context, in *pb.DeleteUserUrlsRequest) (*pb.Empty, error) {
	return nil, nil
}

func (s *Server) GetUrlsByUser(ctx context.Context, in *pb.Empty) (*pb.GetUrlsByUserReply, error) {
	return nil, nil
}

func (s *Server) GetLink(ctx context.Context, in *pb.GetLinkRequest) (*pb.GetLinkReply, error) {
	return nil, nil
}

func (s *Server) SaveLink(ctx context.Context, in *pb.SaveLinkRequest) (*pb.SaveLinkReply, error) {
	return nil, nil
}

func (s *Server) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenReply, error) {
	return nil, nil
}

func (s *Server) ShortenBatch(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchReply, error) {
	return nil, nil
}

func (s *Server) Stats(ctx context.Context, in *pb.Empty) (*pb.StatsReply, error) {
	return nil, nil
}
