package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/handlers"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func setupRouter(s *storage.Storage) *gin.Engine {
	r := gin.Default()
	r.POST("/", handlers.SaveLink(s))
	r.GET("/:id", handlers.GetLink(s))
	return r
}

func NewRouter() *gin.Engine {
	s := storage.New()
	return setupRouter(s)
}
