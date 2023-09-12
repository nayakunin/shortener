package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/config"
)

// AddContext adds config and uuid to gin context
func AddContext(r *gin.Engine, cfg config.Config, userID string) {
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("uuid", userID)
		c.Next()
	})
}
