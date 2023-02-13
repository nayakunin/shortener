package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
)

const paramId = "id"

func GetLinkHandler(s storage.Storager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(paramId)

		link, ok := s.Get(id)
		if !ok {
			// throw error
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, link)
	}
}
