package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func WaitGroup(wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		wg.Add(1)
		c.Next()
		wg.Done()
	}
}
