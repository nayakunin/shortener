package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_Ping(t *testing.T) {
	testCases := []struct {
		name           string
		shortenerError error
		wantErr        bool
		wantReply      *pb.Empty
	}{{
		name:           "success",
		shortenerError: nil,
		wantErr:        false,
		wantReply:      &pb.Empty{},
	}, {
		name:           "error",
		shortenerError: fmt.Errorf("error"),
		wantErr:        true,
		wantReply:      nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortenerService := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				PingError: tc.shortenerError,
			})
			s := NewServer(shortenerService)
			_, err := s.Ping(context.Background(), &pb.Empty{})
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
