package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const paramID = "id"

func (s Server) GetLinkHandler(c *gin.Context) {
	id := c.Param(paramID)

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Empty id"})
		return
	}

	link, ok := s.Storage.Get(id)
	if !ok {
		// throw error
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}
