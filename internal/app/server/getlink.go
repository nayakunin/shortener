package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const paramID = "id"

func (s Server) GetLinkHandler(c *gin.Context) {
	id := c.Param(paramID)

	link, ok := s.storage.Get(id)
	if !ok {
		// throw error
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link)
}
