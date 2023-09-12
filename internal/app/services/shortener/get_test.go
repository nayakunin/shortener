package shortener

import (
	"errors"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_Get(t *testing.T) {
	testCases := []struct {
		name    string
		link    string
		error   error
		wantErr bool
	}{{
		name:    "success",
		link:    "link",
		wantErr: false,
	}, {
		name:    "error",
		error:   errors.New("error"),
		wantErr: true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := testutils.NewMockConfig()
			storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{
				Get: testutils.GetMock{
					Success: tc.link,
					Error:   tc.error,
				},
			})
			s := NewShortenerService(cfg, storage)
			_, err := s.Get("")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
