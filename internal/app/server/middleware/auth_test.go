package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEncodeCookie(t *testing.T) {
	secret := "test"

	testcases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "fGFkNzExNDhjNzlmMjFhYjllZWM1MWVhNWM3ZGQyYjY2ODc5MmY3YzBkMzUzNGFlNjZiMjJmNzFjNjE1MjNmYjM=",
		},
		{
			name:  "non-empty input",
			input: "test",
			want:  "dGVzdHw4OGNkMjEwOGI1MzQ3ZDk3M2NmMzljZGY5MDUzZDdkZDQyNzA0ODc2ZDhjOWE5YmQ4ZTJkMTY4MjU5ZDNkZGY3",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encodeCookie(tc.input, secret)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDecodeCookie(t *testing.T) {
	secret := "test"

	testcases := []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{
			name:  "invalid cookie format",
			input: "",
			want:  "",
			err:   ErrInvalidCookieFormat,
		},
		{
			name:  "non-empty input",
			input: "dGVzdHw4OGNkMjEwOGI1MzQ3ZDk3M2NmMzljZGY5MDUzZDdkZDQyNzA0ODc2ZDhjOWE5YmQ4ZTJkMTY4MjU5ZDNkZGY3",
			want:  "test",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCookie(tc.input, secret)
			if tc.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuth(t *testing.T) {
	secret := "test"

	middleware := Auth(secret)

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(middleware)
	router.GET("/", func(c *gin.Context) {
		userID, ok := c.MustGet("uuid").(string)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusOK, userID)
	})

	t.Run("no cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth",
			Value: "",
		})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("invalid cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth",
			Value: "123",
		})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("valid cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth",
			Value: "dGVzdHw4OGNkMjEwOGI1MzQ3ZDk3M2NmMzljZGY5MDUzZDdkZDQyNzA0ODc2ZDhjOWE5YmQ4ZTJkMTY4MjU5ZDNkZGY3",
		})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test", w.Body.String())
	})
}
