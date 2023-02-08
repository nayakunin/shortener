package getlink

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Storage interface
type Storage interface {
	Get(key string) (string, bool)
}

// Handler handles GET requests
func Handler(s Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		link, ok := s.Get(id)
		if !ok {
			// throw error
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, link)
	}
}
