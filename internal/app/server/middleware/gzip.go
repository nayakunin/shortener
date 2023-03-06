package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
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
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			log.Fatal(c.AbortWithError(500, err))
			return
		}

		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Fatal(c.AbortWithError(500, err))
			}
		}(gz)

		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer = &gzipWriter{c.Writer, gz}
		c.Next()
	}
}
