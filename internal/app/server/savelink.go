package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func (s Server) SaveLinkHandler(c *gin.Context) {
	// read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}(c.Request.Body)

	if len(body) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	urlString := string(body)

	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}

	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// add to storage
	key, err := s.Storage.Add(urlString, userID)
	if err != nil {
		if err == storage.ErrKeyExists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Key already exists"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusCreated, fmt.Sprintf("%s/%s", s.Cfg.BaseURL, key))
}
