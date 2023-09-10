package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
)

// GetUrlsByUserHandler handles GET /urls
func (s Server) GetUrlsByUserHandler(c *gin.Context) {
	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	urls, err := s.Shortener.GetUrlsByUser(userID)
	if err != nil {
		if errors.Is(err, shortener.ErrNoUrlsFound) {
			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "no urls found"})
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, urls)
}
