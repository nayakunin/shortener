package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGzip(b *testing.B) {
	middleware := Gzip()

	inputBytes := []byte("Hello, World!")
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(inputBytes); err != nil {
		b.Fatal(err)
	}
	if err := gzWriter.Close(); err != nil {
		b.Fatal(err)
	}

	body := io.NopCloser(strings.NewReader(buf.String()))

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(middleware)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", body)
	req.Header.Set("Content-Encoding", "gzip")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// send a request
		router.ServeHTTP(w, req)
	}
}

func TestGzip(t *testing.T) {
	middleware := Gzip()

	testString := "Hello, World!"

	inputBytes := []byte(testString)
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(inputBytes); err != nil {
		t.Fatal(err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatal(err)
	}

	body := io.NopCloser(strings.NewReader(buf.String()))
	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.Use(middleware)
	router.GET("/", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatal(err)
		}

		bodyString := string(bodyBytes)
		c.String(http.StatusOK, bodyString)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", body)
	req.Header.Set("Content-Encoding", "gzip")

	router.ServeHTTP(w, req)

	assert.Equal(t, testString, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
