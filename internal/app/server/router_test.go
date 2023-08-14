package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	links := []testutils.MockLink{}
	mockStorage := testutils.NewMockStorage(links)

	cfg := config.Config{}

	router := NewRouter(cfg, mockStorage)

	t.Run("Test Ping Endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
