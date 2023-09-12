package shortener

import (
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_Shorten(t *testing.T) {
	cfg := testutils.NewMockConfig()

	type want struct {
		err  bool
		link string
	}

	type storageReply struct {
		addError   error
		addSuccess string
	}

	testCases := []struct {
		name         string
		input        string
		want         want
		storageReply storageReply
	}{{
		name:  "success",
		input: "https://google.com",
		want: want{
			err:  false,
			link: fmt.Sprintf("%s/%s", cfg.BaseURL, "key"),
		},
		storageReply: storageReply{
			addError:   nil,
			addSuccess: "key",
		},
	}, {
		name:  "error",
		input: "https://google.com",
		want: want{
			err:  true,
			link: fmt.Sprintf("%s/%s", cfg.BaseURL, "key"),
		},
		storageReply: storageReply{
			addError:   fmt.Errorf("error"),
			addSuccess: "key",
		},
	}, {
		name: "invalid url",
		want: want{
			err:  true,
			link: "",
		},
		storageReply: storageReply{
			addError:   fmt.Errorf("error"),
			addSuccess: "key",
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{
				Add: testutils.AddMock{
					Success: tc.storageReply.addSuccess,
					Error:   tc.storageReply.addError,
				},
			})
			s := NewShortenerService(cfg, storage)
			link, err := s.Shorten("user", tc.input)

			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.link, link)
		})
	}
}
