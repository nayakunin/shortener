package shortener

import (
	"errors"
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_DeleteUserUrls(t *testing.T) {
	testCases := []struct {
		name    string
		error   error
		wantErr bool
	}{{
		name:    "success",
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
				DeleteUserUrls: testutils.DeleteUserUrlsMock{
					Error: tc.error,
				},
			})
			s := NewShortenerService(cfg, storage)
			err := s.DeleteUserUrls("", nil)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
