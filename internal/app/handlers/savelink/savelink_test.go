package savelink

import (
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

			body := strings.NewReader(tt.requestBody)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()
			h := Handler(s)
			h.ServeHTTP(w, request)
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
