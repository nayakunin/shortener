// Package server provides server for the application.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/server/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// Server is a struct of the server.
type Server struct {
	Cfg     config.Config
	Storage interfaces.Storage
}

func setupRouter(s Server) (*gin.Engine, *autocert.Manager) {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("config", s.Cfg)
		c.Next()
	})

	r.Use(middleware.Gzip())
	r.Use(middleware.Auth(s.Cfg.AuthSecret))

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

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("nikita-yakunin.dev"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	return r, &m
}

// NewRouter returns a new router for the application
func NewRouter(cfg config.Config, s interfaces.Storage) (*gin.Engine, *autocert.Manager) {
	return setupRouter(Server{
		Cfg:     cfg,
		Storage: s,
	})
}
