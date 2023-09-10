package grpc

import (
	"context"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/storage"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenBatch shortens urls for user
func (s *Server) ShortenBatch(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchReply, error) {
	if in.UserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user_id is required")
	}

	input := make([]interfaces.BatchInput, len(in.Urls))
	for i, url := range in.Urls {
		input[i] = interfaces.BatchInput{
			OriginalURL:   url.OriginalUrl,
			CorrelationID: url.CorrelationId,
		}
	}

	output, err := s.Shortener.ShortenBatch(in.UserID, input)
	if err != nil {
		if errors.Is(err, storage.ErrKeyExists) {
			return nil, status.Error(codes.AlreadyExists, "url already exists")
		}

		if errors.Is(err, storage.ErrBatchInvalidURL) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid url")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	response := make([]*pb.ShortenBatchOutput, len(output))
	for i, v := range output {
		response[i] = &pb.ShortenBatchOutput{
			CorrelationId: v.CorrelationID,
			ShortUrl:      v.ShortURL,
		}
	}

	return &pb.ShortenBatchReply{
		Urls: response,
	}, nil
}
