package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// WaitGroup is a middleware that counts the number of active requests.
func WaitGroup(wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		wg.Add(1)
		c.Next()
		wg.Done()
	}
}
