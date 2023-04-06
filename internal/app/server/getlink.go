package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
	"github.com/pkg/errors"
)

const paramID = "id"

func (s Server) GetLinkHandler(c *gin.Context) {
	id := c.Param(paramID)

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty id"})
		return
	}

	link, err := s.Storage.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrKeyDeleted) {
			c.AbortWithStatusJSON(http.StatusGone, gin.H{"error": "Gone"})
			return
		}

		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}
