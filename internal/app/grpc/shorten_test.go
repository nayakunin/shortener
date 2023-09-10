package grpc

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/storage"
	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_Shorten(t *testing.T) {
	type want struct {
		reply *pb.ShortenReply
		err   bool
	}

	type shortenerReply struct {
		reply string
		err   error
	}

	testCases := []struct {
		name           string
		in             *pb.ShortenRequest
		want           want
		shortenerReply shortenerReply
	}{{
		name: "success",
		in: &pb.ShortenRequest{
			Url:    "https://example.com",
			UserID: "user",
		},
		shortenerReply: shortenerReply{
			reply: "replyUrl",
		},
		want: want{
			reply: &pb.ShortenReply{
				Url: "replyUrl",
			},
		},
	}, {
		name: "error missing user",
		in: &pb.ShortenRequest{
			UserID: "",
			Url:    "https://example.com",
		},
		shortenerReply: shortenerReply{},
		want: want{
			err: true,
		},
	}, {
		name: "error invalid url",
		in: &pb.ShortenRequest{
			UserID: "user",
			Url:    "invalidUrl",
		},
		shortenerReply: shortenerReply{
			err: shortener.ErrInvalidURL,
		},
		want: want{
			err: true,
		},
	}, {
		name: "error already exists",
		in: &pb.ShortenRequest{
			UserID: "user",
			Url:    "https://example.com",
		},
		shortenerReply: shortenerReply{
			err:   storage.ErrKeyExists,
			reply: "replyUrl",
		},
		want: want{
			err: true,
			reply: &pb.ShortenReply{
				Url: "replyUrl",
			},
		},
	}, {
		name: "error internal",
		in: &pb.ShortenRequest{
			UserID: "user",
			Url:    "https://example.com",
		},
		shortenerReply: shortenerReply{
			err: fmt.Errorf("error"),
		},
		want: want{
			err: true,
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortenerService := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				ShortenReplyString: tc.shortenerReply.reply,
				ShortenReplyError:  tc.shortenerReply.err,
			})

			server := NewServer(shortenerService)

			reply, err := server.Shorten(nil, tc.in)
			if tc.want.err {
				assert.Error(t, err)

				if tc.want.reply != nil {
					assert.Equal(t, tc.want.reply.Url, reply.Url)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.reply, reply)
		})
	}
}
