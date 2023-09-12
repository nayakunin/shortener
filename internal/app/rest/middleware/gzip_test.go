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

func TestGzipWriter(t *testing.T) {
	t.Run("Implements", func(t *testing.T) {
		gzWriter := gzipWriter{}
		assert.Implements(t, (*io.Writer)(nil), gzWriter)
		assert.Implements(t, (*io.StringWriter)(nil), gzWriter)
	})

	t.Run("Write", func(t *testing.T) {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		gzw := gzipWriter{writer: writer}

		data := []byte("Hello, world!")
		n, err := gzw.Write(data)

		assert.NoError(t, err)
		assert.Equal(t, len(data), n)

		err = writer.Close()
		assert.NoError(t, err)

		r, err := gzip.NewReader(&buf)
		assert.NoError(t, err)

		decompressedData, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, decompressedData)
	})

	t.Run("WriteString", func(t *testing.T) {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		gzw := gzipWriter{writer: writer}

		data := "Hello, world!"
		n, err := gzw.WriteString(data)

		assert.NoError(t, err)
		assert.Equal(t, len(data), n)

		err = writer.Close()
		assert.NoError(t, err)

		r, err := gzip.NewReader(&buf)
		assert.NoError(t, err)

		decompressedData, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, string(decompressedData))
	})
}

func TestGzip(t *testing.T) {
	type test struct {
		name                     string
		body                     string
		hasContentEncodingHeader bool
		hasAcceptEncodingHeader  bool
		httpCode                 int
		ignoreReply              bool
	}

	tests := []test{{
		name:                     "Content-Encoding header and Accept-Encoding header",
		body:                     compress("Hello, World!"),
		hasContentEncodingHeader: true,
		hasAcceptEncodingHeader:  true,
		httpCode:                 http.StatusOK,
	}, {
		name:                     "No Content-Encoding header and Accept-Encoding header",
		body:                     "Hello, World!",
		hasContentEncodingHeader: false,
		hasAcceptEncodingHeader:  false,
		httpCode:                 http.StatusOK,
	}, {
		name:                     "No Content-Encoding header",
		body:                     "Hello, World!",
		hasContentEncodingHeader: false,
		hasAcceptEncodingHeader:  true,
		httpCode:                 http.StatusOK,
	}, {
		name:                     "No Accept-Encoding header",
		body:                     compress("Hello, World!"),
		hasContentEncodingHeader: true,
		hasAcceptEncodingHeader:  false,
		httpCode:                 http.StatusOK,
	}, {
		name:                     "Empty body",
		body:                     compress(""),
		hasContentEncodingHeader: true,
		hasAcceptEncodingHeader:  true,
		httpCode:                 http.StatusOK,
	}, {
		name:                     "Invalid body",
		body:                     "Hello, World!",
		hasContentEncodingHeader: true,
		hasAcceptEncodingHeader:  true,
		httpCode:                 http.StatusInternalServerError,
		ignoreReply:              true,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tc.body))
			router := gin.New()
			gin.SetMode(gin.TestMode)
			router.Use(Gzip())
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
			if tc.hasContentEncodingHeader {
				req.Header.Set("Content-Encoding", "gzip")
			}

			if tc.hasAcceptEncodingHeader {
				req.Header.Set("Accept-Encoding", "gzip")
			}

			// send a request
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.httpCode, w.Code)

			if tc.ignoreReply {
				return
			}

			testBody := tc.body
			if tc.hasContentEncodingHeader {
				testBody = decompress(tc.body)
			}

			replyBody := w.Body.String()
			if tc.hasAcceptEncodingHeader {
				replyBody = decompress(w.Body.String())
			}

			assert.Equal(t, testBody, replyBody)
		})
	}
}

func compress(data string) string {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	gzw := gzipWriter{writer: writer}

	_, err := gzw.WriteString(data)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func decompress(data string) string {
	r, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		panic(err)
	}

	decompressedData, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(decompressedData)
}
