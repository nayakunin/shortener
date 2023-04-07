package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/server/middleware"
	"github.com/nayakunin/shortener/internal/app/storage"
)

type Server struct {
	Cfg     config.Config
	Storage storage.Storager
}

func setupRouter(s Server) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("config", s.Cfg)
		c.Next()
	})

	r.Use(middleware.Gzip())
	r.Use(middleware.Auth())

	{
		r.POST("/", s.SaveLinkHandler)
		r.GET("/ping", s.pingHandler)
		r.GET("/:id", s.GetLinkHandler)
	}

	api := r.Group("/api")
	{
		api.POST("/shorten", s.ShortenHandler)
		api.POST("/shorten/batch", s.ShortenBatchHandler)
		api.GET("/user/urls", s.GetUrlsByUserHandler)
		api.DELETE("/user/urls", s.DeleteUserUrlsHandler)
	}
	return r
}

func NewRouter(cfg config.Config, s storage.Storager) *gin.Engine {
	return setupRouter(Server{
		Cfg:     cfg,
		Storage: s,
	})
}
