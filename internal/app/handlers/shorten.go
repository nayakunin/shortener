package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func ShortenHandler(s storage.Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var req ShortenRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
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

		key, err := s.Add(req.URL)
		if err != nil {
			fmt.Println(err)
			if err == storage.ErrKeyExists {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Key already exists"})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, ShortenResponse{Result: fmt.Sprintf("%s/%s", config.Config.BaseURL, key)})
	}
}
