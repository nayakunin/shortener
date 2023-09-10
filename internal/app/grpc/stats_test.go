package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_Stats(t *testing.T) {
	type want struct {
		reply *pb.StatsReply
		err   bool
	}

	type shortenerReply struct {
		reply *interfaces.Stats
		err   error
	}

	testCases := []struct {
		name           string
		want           want
		shortenerReply shortenerReply
	}{{
		name: "success",
		want: want{
			reply: &pb.StatsReply{
				Urls:  1,
				Users: 2,
			},
		},
		shortenerReply: shortenerReply{
			reply: &interfaces.Stats{
				Urls:  1,
				Users: 2,
			},
		},
	}, {
		name: "error",
		want: want{
			err: true,
		},
		shortenerReply: shortenerReply{
			err: fmt.Errorf("error"),
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				StatsReply: tc.shortenerReply.reply,
				StatsError: tc.shortenerReply.err,
			})
			s := NewServer(store)
			reply, err := s.Stats(context.Background(), &pb.Empty{})

			if tc.want.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.reply, reply)
		})
	}
}
