package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetUrlsByUser(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name                string
		request             string
		links               []testutils.MockLink
		shouldCheckLocation bool
		want                want
	}{
		{
			name:    "not found",
			request: "/api/user/urls",
			want: want{
				statusCode:  http.StatusNoContent,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:    "success",
			request: "/api/user/urls",
			links: []testutils.MockLink{
				{
					OriginalURL: "https://google.com",
					ShortURL:    "link",
					UserID:      "userID",
				},
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json; charset=utf-8",
				body:        "[{\"original_url\":\"https://google.com\",\"short_url\":\"http://localhost:8080/link\"}]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutils.NewMockStorage(tt.links)
			cfg := testutils.NewMockConfig()
			server := Server{
				Storage: s,
				Cfg:     cfg,
			}

			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			router.GET("/api/user/urls", server.GetUrlsByUserHandler)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			if tt.shouldCheckLocation {
				var link testutils.MockLink
				for _, l := range tt.links {
					if l.ShortURL == tt.request[1:] {
						link = l
						break
					}
				}
				assert.Equal(t, link, res.Header.Get("Location"))
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			if tt.want.body != "" {
				assert.JSONEq(t, tt.want.body, w.Body.String())
			}
		})
	}
}
