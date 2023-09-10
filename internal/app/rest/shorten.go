package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/services/shortener"
	"github.com/nayakunin/shortener/internal/app/storage"
)

// ShortenRequest request structure for POST /shorten
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse response structure for POST /shorten
type ShortenResponse struct {
	Result string `json:"result"`
}

// ShortenHandler handles POST /shorten
// Body ShortenRequest
func (s Server) ShortenHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	var req ShortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.URL) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	shortUrl, err := s.Shortener.Shorten(userID, req.URL)
	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
			return
		}

		if errors.Is(err, storage.ErrKeyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, ShortenResponse{Result: shortUrl})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal rest error"})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{Result: shortUrl})
}
