package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
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
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "no urls found"})
		return
	}

	response := make([]storage.Link, 0, len(urls))
	for shortURL, originalURL := range urls {
		response = append(response, storage.Link{
			ShortURL:    fmt.Sprintf("%s/%s", s.Cfg.BaseURL, shortURL),
			OriginalURL: originalURL,
		})
	}

	c.JSON(http.StatusOK, response)
}
