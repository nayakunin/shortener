package shortener

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_Stats(t *testing.T) {
	type want struct {
		err   bool
		stats *interfaces.Stats
	}

	type storageReply struct {
		statsError   error
		statsSuccess interfaces.Stats
	}

	testCases := []struct {
		name         string
		want         want
		storageReply storageReply
	}{{
		name: "success",
		want: want{
			err: false,
			stats: &interfaces.Stats{
				Urls:  1,
				Users: 1,
			},
		},
		storageReply: storageReply{
			statsError: nil,
			statsSuccess: interfaces.Stats{
				Urls:  1,
				Users: 1,
			},
		},
	}, {
		name: "error",
		want: want{
			err: true,
		},
		storageReply: storageReply{
			statsError: fmt.Errorf("error"),
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := testutils.NewMockConfig()
			storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{
				Stats: testutils.StatsMock{
					Error:   tc.storageReply.statsError,
					Success: tc.storageReply.statsSuccess,
				},
			})
			s := NewShortenerService(cfg, storage)
			stats, err := s.Stats()

			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want.stats, stats)
		})
	}
}
