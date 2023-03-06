package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/handlers"
	"github.com/nayakunin/shortener/internal/app/server/middleware"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func setupRouter(s storage.Storager) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Gzip())
	r.POST("/", handlers.SaveLinkHandler(s))
	r.GET("/:id", handlers.GetLinkHandler(s))
	r.POST("/api/shorten", handlers.ShortenHandler(s))
	return r
}

func NewRouter(s *storage.Storage) *gin.Engine {
	return setupRouter(s)
}
