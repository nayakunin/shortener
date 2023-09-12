package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	type want struct {
		statusCode int
		reply      *interfaces.Stats
	}

	type shortenerReply struct {
		reply *interfaces.Stats
		err   error
	}

	tests := []struct {
		name           string
		want           want
		shortenerReply shortenerReply
	}{
		{
			name: "success",
			want: want{
				statusCode: http.StatusOK,
				reply: &interfaces.Stats{
					Urls:  1,
					Users: 2,
				},
			},
			shortenerReply: shortenerReply{
				reply: &interfaces.Stats{
					Urls:  1,
					Users: 2,
				},
			},
		},
		{
			name: "error",
			want: want{
				statusCode: http.StatusInternalServerError,
			},
			shortenerReply: shortenerReply{
				err: fmt.Errorf("error"),
			},
		},
	}

	cfg := testutils.NewMockConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testutils.NewMockShortenerService(testutils.MockShortenerServiceParameters{
				StatsReply: tt.shortenerReply.reply,
				StatsError: tt.shortenerReply.err,
			})
			router := gin.Default()
			testutils.AddContext(router, cfg, "userID")
			server := Server{
				Shortener: service,
			}
			router.GET("/", server.statsHandler)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			if tt.want.reply != nil {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				replyString := fmt.Sprintf(`{"urls":%d,"users":%d}`, tt.want.reply.Urls, tt.want.reply.Users)
				assert.Equal(t, replyString, string(resBody))
			}
		})
	}
}
