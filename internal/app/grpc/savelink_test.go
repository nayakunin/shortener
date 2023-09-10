package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/storage"
	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_SaveLink(t *testing.T) {
	type want struct {
		reply *pb.SaveLinkReply
		err   bool
	}

	type shortenerReply struct {
		reply string
		err   error
	}

	testCases := []struct {
		name           string
		in             *pb.SaveLinkRequest
		want           want
		shortenerReply shortenerReply
	}{{
		name: "success",
		in: &pb.SaveLinkRequest{
			Url:    "https://example.com",
			UserID: "user",
		},
		shortenerReply: shortenerReply{
			reply: "replyUrl",
		},
		want: want{
			reply: &pb.SaveLinkReply{
				Url: "replyUrl",
			},
		},
	}, {
		name: "error missing user",
		in: &pb.SaveLinkRequest{
			UserID: "",
			Url:    "https://example.com",
		},
		shortenerReply: shortenerReply{},
		want: want{
			err: true,
		},
	}, {
		name: "error invalid url",
		in: &pb.SaveLinkRequest{
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
		in: &pb.SaveLinkRequest{
			UserID: "user",
			Url:    "https://example.com",
		},
		shortenerReply: shortenerReply{
			err:   storage.ErrKeyExists,
			reply: "replyUrl",
		},
		want: want{
			err: true,
			reply: &pb.SaveLinkReply{
				Url: "replyUrl",
			},
		},
	}, {
		name: "error internal",
		in: &pb.SaveLinkRequest{
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

			reply, err := server.SaveLink(context.Background(), tc.in)
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
