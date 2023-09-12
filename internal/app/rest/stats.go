package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s Server) statsHandler(c *gin.Context) {
	stats, err := s.Shortener.Stats()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, stats)
}
