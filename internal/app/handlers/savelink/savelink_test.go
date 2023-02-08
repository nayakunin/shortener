package savelink

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockStorage struct {
	links map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		links: make(map[string]string),
	}
}

func (m *mockStorage) Add(link string) string {
	key := "link"
	m.links[key] = link
	return key
}

func setupRouter(s *mockStorage) *gin.Engine {
	r := gin.Default()
	r.POST("/", Handler(s))
	return r
}

func TestSaveLink(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		response    string
	}

	tests := []struct {
		name                string
		requestBody         string
		shouldCheckResponse bool
		want                want
	}{
		{
			name:        "success",
			requestBody: "https://google.com",
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/link",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockStorage()

			router := setupRouter(s)
			w := httptest.NewRecorder()
			body := strings.NewReader(tt.requestBody)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			router.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
