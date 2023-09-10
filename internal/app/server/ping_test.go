package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("Regular storage", func(t *testing.T) {
		s := testutils.NewMockStorage([]testutils.MockLink{})
		router := gin.Default()
		cfg := testutils.NewMockConfig()
		service := shortener.NewShortenerService(cfg, s)
		server := Server{
			Shortener: service,
		}
		router.GET("/ping", server.pingHandler)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
