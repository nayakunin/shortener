package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShorten(t *testing.T) {
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
			requestBody:         `{"url": "https://google.com"}`,
			shouldCheckResponse: true,
			want: want{
				statusCode:  http.StatusCreated,
				response:    fmt.Sprintf(`{"result":"%s/%s"}`, cfg.BaseURL, "link"),
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "empty body",
			requestBody: `{ "url": "" }`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "invalid url",
			requestBody: `{"url": "google.com"}`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "duplicate url",
			requestBody: `{"url": "https://google.com"}`,
			links: []testutils.MockLink{
				{
					OriginalURL: "https://google.com",
					ShortURL:    "link",
				},
			},
			shouldCheckResponse: true,
			want: want{
				statusCode:  http.StatusConflict,
				response:    fmt.Sprintf(`{"result":"%s/%s"}`, cfg.BaseURL, "link"),
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutils.NewMockStorage(tt.links)
			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			server := Server{
				Storage: s,
				Cfg:     cfg,
			}
			router.POST("/", server.ShortenHandler)

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
