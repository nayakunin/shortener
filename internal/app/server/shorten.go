package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (s Server) ShortenHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}

	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	key, err := s.Storage.Add(req.URL, userID)
	if err != nil {
		fmt.Println(err)
		if err == storage.ErrKeyExists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Key already exists"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{Result: fmt.Sprintf("%s/%s", s.Cfg.BaseURL, key)})
}
