package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s Server) pingHandler(c *gin.Context) {
	err := s.Shortener.Ping()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
