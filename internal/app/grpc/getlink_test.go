package grpc

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetLink(t *testing.T) {
	type want struct {
		reply *pb.GetLinkReply
		err   bool
	}

	type StorageReply struct {
		GetSuccess string
		GetError   error
	}

	testCases := []struct {
		name         string
		in           *pb.GetLinkRequest
		want         want
		storageReply StorageReply
	}{{
		name: "success",
		in: &pb.GetLinkRequest{
			Key: "key",
		},
		want: want{
			reply: &pb.GetLinkReply{
				Url: "url",
			},
		},
		storageReply: StorageReply{
			GetSuccess: "url",
		},
	}, {
		name: "error",
		in: &pb.GetLinkRequest{
			Key: "key",
		},
		want: want{
			err: true,
		},
		storageReply: StorageReply{
			GetError: fmt.Errorf("error"),
		},
	}, {
		name: "empty key",
		in: &pb.GetLinkRequest{
			Key: "",
		},
		want: want{
			err: true,
		},
		storageReply: StorageReply{
			GetError: fmt.Errorf("error"),
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortenerService := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				GetReplyError:  tc.storageReply.GetError,
				GetReplyString: tc.storageReply.GetSuccess,
			})
			s := NewServer(shortenerService)
			reply, err := s.GetLink(nil, tc.in)

			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.reply, reply)
		})
	}
}
