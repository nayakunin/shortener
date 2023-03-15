package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("auth")
		if err != nil {
			if err == http.ErrNoCookie {
				c.SetCookie("auth", uuid.NewString(), 3600, "/", "", false, true)
				c.Next()
				return
			}

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Next()
	}
}
