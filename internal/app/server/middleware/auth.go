package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var secret = []byte("secret")

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		encryptedCookie, cookieReadingErr := c.Cookie("auth")
		if cookieReadingErr != nil {
			if cookieReadingErr == http.ErrNoCookie {
				userId := uuid.NewString()
				h := hmac.New(sha256.New, secret)
				h.Write([]byte(userId))
				cookie := h.Sum(nil)

				c.SetCookie("auth", string(cookie), 3600, "/", "", false, true)
				c.Set("uuid", userId)
				c.Next()
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		decryptedCookie, err := hex.DecodeString(encryptedCookie)
		if err != nil {
			return
		}

		h := hmac.New(sha256.New, secret)
		h.Write(decryptedCookie)
		sign := h.Sum(nil)

		if !hmac.Equal(sign, decryptedCookie) {
			userId := uuid.NewString()
			h := hmac.New(sha256.New, secret)
			h.Write([]byte(userId))
			cookie := h.Sum(nil)

			c.SetCookie("auth", string(cookie), 3600, "/", "", false, true)
			c.Set("uuid", userId)
			c.Next()
			return
		}

		c.Next()
	}
}
