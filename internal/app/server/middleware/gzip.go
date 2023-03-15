package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w gzipWriter) WriteString(s string) (int, error) {
	return w.writer.Write([]byte(s))
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				log.Print(c.AbortWithError(http.StatusInternalServerError, err))
				return
			}

			defer func(gzr *gzip.Reader) {
				err := gzr.Close()
				if err != nil {
					log.Print(c.AbortWithError(http.StatusInternalServerError, err))
				}
			}(reader)

			body, err := io.ReadAll(reader)
			if err != nil {
				log.Print(c.AbortWithError(http.StatusInternalServerError, err))
				return
			}

			c.Request.Body = io.NopCloser(strings.NewReader(string(body)))
			c.Request.Header.Del("Content-Encoding")
			c.Request.ContentLength = int64(len(body))
		}

		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			log.Print(c.AbortWithError(http.StatusInternalServerError, err))
			return
		}

		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Print(c.AbortWithError(http.StatusInternalServerError, err))
			}
		}(gz)

		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer = &gzipWriter{c.Writer, gz}
		c.Next()
	}
}
