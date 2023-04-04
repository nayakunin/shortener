package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserUrls(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}

	tests := []struct {
		name                string
		request             string
		links               []testutils.MockLink
		shouldCheckLocation bool
		want                want
	}{
		{
			name:    "return accepted",
			request: "[]",
			want: want{
				statusCode:  http.StatusAccepted,
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:    "empty body",
			request: "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	cfg := testutils.NewMockConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutils.NewMockStorage(tt.links)
			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			server := Server{
				Storage: s,
				Cfg:     cfg,
			}
			router.GET("/", server.DeleteUserUrlsHandler)

			w := httptest.NewRecorder()
			body := strings.NewReader(tt.request)
			request := httptest.NewRequest(http.MethodGet, "/", body)
			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
