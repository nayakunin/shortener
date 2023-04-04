package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func (s Server) pingHandler(c *gin.Context) {
	dbStorage, ok := s.Storage.(*storage.DBStorage)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	}

	if err := dbStorage.Connection.Ping(context.Background()); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
