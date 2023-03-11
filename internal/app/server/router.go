package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/handlers"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/server/middleware"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func setupRouter(cfg config.Config, s storage.Storager) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	r.Use(middleware.Gzip())
	r.POST("/", handlers.SaveLinkHandler(s))
	r.GET("/:id", handlers.GetLinkHandler(s))
	r.POST("/api/shorten", handlers.ShortenHandler(s))
	return r
}

func NewRouter(cfg config.Config, s storage.Storager) *gin.Engine {
	return setupRouter(cfg, s)
}
