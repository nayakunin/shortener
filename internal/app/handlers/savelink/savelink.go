package savelink

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// Storage interface
type Storage interface {
	Add(link string) string
}

// Handler handles POST requests
func Handler(s Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// read body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}(c.Request.Body)

		url := string(body)

		// add to storage
		key := s.Add(url)

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusCreated, fmt.Sprintf("http://localhost:8080/%s", key))
	}

}
