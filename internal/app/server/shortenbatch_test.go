package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenBatch(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		response    []ShortenBatchOutput
	}

	cfg := testutils.NewMockConfig()

	tests := []struct {
		name                string
		requestBody         []ShortenBatchInput
		shouldCheckResponse bool
		storage             []testutils.MockLink
		want                want
	}{
		{
			name: "success",
			requestBody: []ShortenBatchInput{{
				CorrelationID: "1",
				OriginalURL:   "https://google.com/1",
			}},
			shouldCheckResponse: true,
			want: want{
				statusCode: http.StatusCreated,
				response: []ShortenBatchOutput{{
					CorrelationID: "1",
					ShortURL:      fmt.Sprintf("%s/%s", cfg.BaseURL, "link"),
				}},
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:        "empty body",
			requestBody: []ShortenBatchInput{},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "invalid url",
			requestBody: []ShortenBatchInput{{
				CorrelationID: "1",
				OriginalURL:   "sda",
			}},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "duplicate url",
			requestBody: []ShortenBatchInput{{
				CorrelationID: "1",
				OriginalURL:   "https://google.com",
			}},
			storage: []testutils.MockLink{
				{
					OriginalURL: "https://google.com",
					ShortURL:    "link",
				},
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutils.NewMockStorage(tt.storage)
			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			service := shortener.NewShortenerService(cfg, s)
			server := Server{
				Shortener: service,
			}
			router.POST("/", server.ShortenBatchHandler)

			w := httptest.NewRecorder()
			input, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			body := strings.NewReader(string(input))
			request := httptest.NewRequest(http.MethodPost, "/", body)
			router.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()

			if tt.shouldCheckResponse {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				output, err := json.Marshal(tt.want.response)
				require.NoError(t, err)
				assert.JSONEq(t, string(output), string(resBody))
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
