package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/middleware"
)

func setupRouter(s Server) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("config", s.cfg)
		c.Next()
	})

	r.Use(middleware.Gzip())
	r.POST("/", s.SaveLinkHandler)
	r.GET("/:id", s.GetLinkHandler)
	r.POST("/api/shorten", s.ShortenHandler)
	return r
}

func NewRouter(s Server) *gin.Engine {
	return setupRouter(s)
}
