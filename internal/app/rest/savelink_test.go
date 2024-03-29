package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLink(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		response    string
	}

	cfg := testutils.NewMockConfig()

	tests := []struct {
		name                string
		requestBody         string
		shouldCheckResponse bool
		links               []testutils.MockLink
		want                want
	}{
		{
			name:                "success",
			requestBody:         "https://google.com",
			shouldCheckResponse: true,
			want: want{
				statusCode:  http.StatusCreated,
				response:    fmt.Sprintf("%s/%s", cfg.BaseURL, "link"),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "empty body",
			requestBody: "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid url",
			requestBody: "google.com",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "duplicate url",
			requestBody: "https://google.com",
			links: []testutils.MockLink{
				{
					OriginalURL: "https://google.com",
					ShortURL:    "link",
				},
			},
			shouldCheckResponse: true,
			want: want{
				statusCode:  http.StatusConflict,
				response:    fmt.Sprintf("%s/%s", cfg.BaseURL, "link"),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutils.NewMockStorage(tt.links)
			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			service := shortener.NewShortenerService(cfg, s)
			server := Server{
				Shortener: service,
			}
			router.POST("/", server.SaveLinkHandler)

			w := httptest.NewRecorder()
			body := strings.NewReader(tt.requestBody)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			router.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()

			if tt.shouldCheckResponse {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.want.response, string(resBody))
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
