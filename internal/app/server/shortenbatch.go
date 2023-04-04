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

type ShortenBatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s Server) ShortenBatchHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var req []ShortenBatchInput
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	input := make([]storage.BatchInput, len(req))
	for i, v := range req {
		_, err = url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
			return
		}

		input[i] = storage.BatchInput{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
		}
	}

	output, err := s.Storage.AddBatch(input, userID)
	if err != nil {
		if err == storage.ErrKeyExists {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Key already exists"})
			return
		}

		if err == storage.ErrBatchInvalidURL {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := make([]ShortenBatchOutput, len(output))
	for i, v := range output {
		response[i] = ShortenBatchOutput{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.Cfg.BaseURL, v.Key),
		}
	}

	c.JSON(http.StatusCreated, response)
}
