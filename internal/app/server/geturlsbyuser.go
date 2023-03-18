package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
	"net/http"
)

func (s Server) GetUrlsByUserHandler(c *gin.Context) {
	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	urls, err := s.Storage.GetUrlsByUser(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": "No urls found"})
		return
	}

	response := make([]storage.Link, 0, len(urls))
	for shortLink, originalLink := range urls {
		response = append(response, storage.Link{
			ShortURL:    shortLink,
			OriginalURL: originalLink,
		})
	}

	c.JSON(http.StatusOK, response)
}
