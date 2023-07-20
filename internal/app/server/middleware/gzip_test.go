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
	gin.SetMode(gin.ReleaseMode)
	router.Use(middleware)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	w := httptest.NewRecorder()
	w.Header().Set("Content-Encoding", "gzip")
	req, _ := http.NewRequest("GET", "/", body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// send a request
		router.ServeHTTP(w, req)
	}
}