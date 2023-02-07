package getlink

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStorage struct {
	links map[string]string
}

func newMockStorage(links map[string]string) *mockStorage {
	return &mockStorage{
		links: links,
	}
}

func (s *mockStorage) Get(key string) (string, bool) {
	link, ok := s.links[key]
	return link, ok
}

func TestGetLink(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}

	tests := []struct {
		name                string
		request             string
		links               map[string]string
		shouldCheckLocation bool
		want                want
	}{
		{
			name:    "bad request",
			request: "/",
			links:   map[string]string{},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "not found",
			request: "/link",
			links:   map[string]string{},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "success",
			request: "/link",
			links: map[string]string{
				"link": "https://google.com",
			},
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockStorage(tt.links)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := Handler(s)
			h.ServeHTTP(w, request)
			res := w.Result()

			if tt.shouldCheckLocation {
				assert.Equal(t, tt.links["link"], res.Header.Get("Location"))
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
