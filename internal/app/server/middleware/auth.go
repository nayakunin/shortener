package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalidCookieFormat    = errors.New("invalid cookie format")
	ErrInvalidCookieSignature = errors.New("invalid cookie signature")
)

var secret = "secret"

func encodeCookie(value, secretKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(value))
	signature := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(value + "|" + signature)), nil
}

func decodeCookie(encoded, secretKey string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return "", ErrInvalidCookieFormat
	}

	value, signature := parts[0], parts[1]

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(value))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", ErrInvalidCookieSignature
	}

	return value, nil
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		var value string
		if err != nil || cookie == "" {
			value = uuid.NewString()
			encoded, err := encodeCookie(value, secret)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encoded, 3600, "/", "", false, true)
			c.Set("uuid", value)
			c.Next()
		}

		value, err = decodeCookie(cookie, secret)
		if err != nil {
			value = uuid.NewString()
			encoded, err := encodeCookie(value, secret)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie("auth", encoded, 3600, "/", "", false, true)
		}

		c.Set("uuid", value)
		c.Next()
	}
}
