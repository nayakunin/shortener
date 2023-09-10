package grpc

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_DeleteUserUrls(t *testing.T) {
	type shortenerReply struct {
		err error
	}

	testCases := []struct {
		name           string
		in             *pb.DeleteUserUrlsRequest
		wantErr        bool
		shortenerReply shortenerReply
	}{{
		name: "success",
		in: &pb.DeleteUserUrlsRequest{
			UserID: "user",
			Keys:   []string{"key"},
		},
		wantErr: false,
		shortenerReply: shortenerReply{
			err: nil,
		},
	}, {
		name: "error",
		in: &pb.DeleteUserUrlsRequest{
			UserID: "user",
			Keys:   []string{"key"},
		},
		wantErr: true,
		shortenerReply: shortenerReply{
			err: fmt.Errorf("error"),
		},
	}, {
		name: "empty user id",
		in: &pb.DeleteUserUrlsRequest{
			UserID: "",
			Keys:   []string{"key"},
		},
		wantErr: true,
		shortenerReply: shortenerReply{
			err: nil,
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortenerService := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				DeleteReplyError: tc.shortenerReply.err,
			})
			s := NewServer(shortenerService)
			_, err := s.DeleteUserUrls(nil, tc.in)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
