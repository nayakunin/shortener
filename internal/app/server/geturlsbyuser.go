package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s Server) GetUrlsByUserHandler(c *gin.Context) {
	userId, ok := c.MustGet("uuid").(string)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	urls, err := s.Storage.GetUrlsByUser(userId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, urls)
}
