package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteUserUrlsHandler deletes urls from user's storage.
func (s Server) DeleteUserUrlsHandler(c *gin.Context) {
	userID, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get keys array from body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var keys []string
	err = json.Unmarshal(body, &keys)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := s.Storage.DeleteUserUrls(userID, keys); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
}
