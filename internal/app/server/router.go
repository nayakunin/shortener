package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/handlers"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func setupRouter(s storage.Storager) *gin.Engine {
	r := gin.Default()
	r.POST("/", handlers.SaveLinkHandler(s))
	r.GET("/:id", handlers.GetLinkHandler(s))
	r.POST("/api/shorten", handlers.ShortenHandler(s))
	return r
}

func NewRouter() *gin.Engine {
	s := storage.New()
	return setupRouter(s)
}
