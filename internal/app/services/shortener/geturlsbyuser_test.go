package shortener

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUrlsByUser(t *testing.T) {
	cfg := testutils.NewMockConfig()

	testCases := []struct {
		name    string
		links   map[string]string
		error   error
		wantErr bool
		want    []interfaces.Link
	}{{
		name:  "success",
		links: map[string]string{"key": "link"},
		want: []interfaces.Link{{
			ShortURL:    fmt.Sprintf("%s/%s", cfg.BaseURL, "key"),
			OriginalURL: "link",
		}},
		wantErr: false,
	}, {
		name:    "no urls found",
		links:   map[string]string{},
		wantErr: true,
	}, {
		name:    "error",
		error:   errors.New("error"),
		wantErr: true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{
				GetUrlsByUser: testutils.GetUrlsByUserMock{
					Success: tc.links,
					Error:   tc.error,
				},
			})
			s := NewShortenerService(cfg, storage)
			_, err := s.GetUrlsByUser("")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
