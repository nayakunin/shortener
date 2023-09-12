package rest

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/storage"
)

// SaveLinkHandler handles POST /
// Request body contains url to save
func (s Server) SaveLinkHandler(c *gin.Context) {
	// read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
			return
		}
	}(c.Request.Body)

	if len(body) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	urlString := string(body)

	// add to storage
	shortURL, err := s.Shortener.Shorten(userID, urlString)
	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
			return
		}

		if errors.Is(err, storage.ErrKeyExists) {
			c.String(http.StatusConflict, shortURL)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusCreated, shortURL)
}
