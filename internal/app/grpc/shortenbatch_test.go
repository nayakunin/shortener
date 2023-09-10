package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/storage"
	"github.com/nayakunin/shortener/internal/app/testutils"
	pb "github.com/nayakunin/shortener/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_ShortenBatch(t *testing.T) {
	type want struct {
		reply *pb.ShortenBatchReply
		err   bool
	}

	type shortenerReply struct {
		reply []interfaces.BatchOutput
		err   error
	}

	testCases := []struct {
		name           string
		in             *pb.ShortenBatchRequest
		want           want
		shortenerReply shortenerReply
	}{{
		name: "success",
		in: &pb.ShortenBatchRequest{
			UserID: "user",
			Urls: []*pb.ShortenBatchInput{{
				CorrelationId: "1",
				OriginalUrl:   "https://example.com",
			}, {
				CorrelationId: "2",
				OriginalUrl:   "https://example.com/2",
			}},
		},
		shortenerReply: shortenerReply{
			reply: []interfaces.BatchOutput{{
				ShortURL:      "short1",
				CorrelationID: "1",
			}, {
				ShortURL:      "short2",
				CorrelationID: "2",
			}},
		},
		want: want{
			reply: &pb.ShortenBatchReply{
				Urls: []*pb.ShortenBatchOutput{{
					CorrelationId: "1",
					ShortUrl:      "short1",
				}, {
					CorrelationId: "2",
					ShortUrl:      "short2",
				}},
			},
			err: false,
		},
	}, {
		name: "error missing user",
		in: &pb.ShortenBatchRequest{
			UserID: "",
		},
		want: want{
			err: true,
		},
	}, {
		name: "error duplicate url",
		in: &pb.ShortenBatchRequest{
			UserID: "user",
			Urls: []*pb.ShortenBatchInput{{
				CorrelationId: "1",
				OriginalUrl:   "https://example.com",
			}},
		},
		shortenerReply: shortenerReply{
			err: storage.ErrKeyExists,
		},
		want: want{
			err: true,
		},
	}, {
		name: "error invalid url",
		in: &pb.ShortenBatchRequest{
			UserID: "user",
			Urls: []*pb.ShortenBatchInput{{
				CorrelationId: "1",
				OriginalUrl:   "invalidUrl",
			}},
		},
		shortenerReply: shortenerReply{
			err: storage.ErrBatchInvalidURL,
		},
		want: want{
			err: true,
		},
	}, {
		name: "error internal",
		in: &pb.ShortenBatchRequest{
			UserID: "user",
			Urls: []*pb.ShortenBatchInput{{
				CorrelationId: "1",
				OriginalUrl:   "https://example.com",
			}},
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
			store := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				ShortenBatchReply: tc.shortenerReply.reply,
				ShortenBatchError: tc.shortenerReply.err,
			})
			s := NewServer(store)
			reply, err := s.ShortenBatch(context.Background(), tc.in)

			if tc.want.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.reply, reply)
		})
	}
}
