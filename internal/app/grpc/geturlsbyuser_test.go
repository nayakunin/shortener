package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetUrlsByUser(t *testing.T) {
	type want struct {
		reply *pb.GetUrlsByUserReply
		err   bool
	}

	type shortenerReply struct {
		reply []interfaces.Link
		err   error
	}

	testCases := []struct {
		name           string
		in             *pb.GetUrlsByUserRequest
		want           want
		shortenerReply shortenerReply
	}{{
		name: "success",
		in: &pb.GetUrlsByUserRequest{
			UserID: "user",
		},
		want: want{
			reply: &pb.GetUrlsByUserReply{
				Urls: []*pb.GetUrlsByUserLink{{
					ShortUrl:    "short",
					OriginalUrl: "original",
				}},
			},
			err: false,
		},
		shortenerReply: shortenerReply{
			err: nil,
			reply: []interfaces.Link{{
				ShortURL:    "short",
				OriginalURL: "original",
			}},
		},
	}, {
		name: "error",
		in: &pb.GetUrlsByUserRequest{
			UserID: "user",
		},
		want: want{
			err: true,
		},
		shortenerReply: shortenerReply{
			err: fmt.Errorf("error"),
		},
	}, {
		name: "empty user id",
		in: &pb.GetUrlsByUserRequest{
			UserID: "",
		},
		want: want{
			err: true,
		},
		shortenerReply: shortenerReply{
			err: nil,
		},
	}, {
		name: "empty reply from the service",
		in: &pb.GetUrlsByUserRequest{
			UserID: "user",
		},
		want: want{
			err: true,
		},
		shortenerReply: shortenerReply{
			err: shortener.ErrNoUrlsFound,
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortenerService := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				GetUrlsReply: tc.shortenerReply.reply,
				GetUrlsError: tc.shortenerReply.err,
			})
			s := NewServer(shortenerService)
			_, err := s.GetUrlsByUser(context.Background(), tc.in)

			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
