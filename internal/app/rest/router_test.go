package rest

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/nayakunin/shortener/internal/app/config"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	links := []testutils.MockLink{}
	mockStorage := testutils.NewMockStorage(links)

	cfg := config.Config{}

	wg := &sync.WaitGroup{}
	service := shortener.NewShortenerService(cfg, mockStorage)

	router, _ := NewRouter(service, wg, cfg.AuthSecret, cfg.TrustedSubnet)

	t.Run("Test Ping Endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
