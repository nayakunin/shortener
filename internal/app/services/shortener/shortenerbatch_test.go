package shortener

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_ShortenBatch(t *testing.T) {
	cfg := testutils.NewMockConfig()

	type want struct {
		err   bool
		links []interfaces.BatchOutput
	}

	type storageReply struct {
		addError   error
		addSuccess []interfaces.DBBatchOutput
	}

	testCases := []struct {
		name         string
		input        []interfaces.BatchInput
		want         want
		storageReply storageReply
	}{{
		name: "success",
		input: []interfaces.BatchInput{{
			OriginalURL:   "https://google.com",
			CorrelationID: "1",
		}, {
			OriginalURL:   "https://google.com/2",
			CorrelationID: "2",
		}},
		storageReply: storageReply{
			addError: nil,
			addSuccess: []interfaces.DBBatchOutput{{
				Key:           "key1",
				CorrelationID: "1",
			}, {
				Key:           "key2",
				CorrelationID: "2",
			}},
		},
		want: want{
			err: false,
			links: []interfaces.BatchOutput{{
				ShortURL:      fmt.Sprintf("%s/%s", cfg.BaseURL, "key1"),
				CorrelationID: "1",
			}, {
				ShortURL:      fmt.Sprintf("%s/%s", cfg.BaseURL, "key2"),
				CorrelationID: "2",
			}},
		},
	}, {
		name:  "error",
		input: []interfaces.BatchInput{},
		storageReply: storageReply{
			addError: fmt.Errorf("error"),
		},
		want: want{
			err: true,
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{
				AddBatch: testutils.AddBatchMock{
					Success: tc.storageReply.addSuccess,
					Error:   tc.storageReply.addError,
				},
			})
			s := NewShortenerService(cfg, storage)
			links, err := s.ShortenBatch("userID", tc.input)

			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.links, links)
		})
	}
}
