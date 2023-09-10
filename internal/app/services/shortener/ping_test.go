package shortener

import (
	"testing"

	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestService_Ping(t *testing.T) {
	t.Run("Regular storage", func(t *testing.T) {
		storage := testutils.NewSimpleMockStorage(testutils.SimpleMockStorageParameters{})

		s := NewShortenerService(testutils.NewMockConfig(), storage)
		err := s.Ping()

		assert.NoError(t, err)
	})
}
