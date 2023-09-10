// Package rest provides rest for the application.
package rest

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/rest/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type Shortener interface {
	Shorten(userID string, url string) (string, error)
	ShortenBatch(userID string, urls []interfaces.BatchInput) ([]interfaces.BatchOutput, error)
	Get(key string) (string, error)
	DeleteUserUrls(userID string, keys []string) error
	GetUrlsByUser(userID string) ([]interfaces.Link, error)
	Stats() (*interfaces.Stats, error)
	Ping() error
}

// Server is a struct of the rest.
type Server struct {
	Shortener Shortener
}

func setupRouter(wg *sync.WaitGroup, s Server, authSecret string, trustedSubnet string) (*gin.Engine, *autocert.Manager) {
	r := gin.Default()

	r.Use(middleware.Gzip())
	r.Use(middleware.Auth(authSecret))
	r.Use(middleware.WaitGroup(wg))

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

		internal := api.Group("/internal")
		internal.Use(middleware.Internal(trustedSubnet))
		{
			internal.GET("/internal/stats", s.statsHandler)
		}
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("nikita-yakunin.dev"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	return r, &m
}

// NewRouter returns a new router for the application
func NewRouter(shortener Shortener, wg *sync.WaitGroup, authSecret string, trustedSubnet string) (*gin.Engine, *autocert.Manager) {
	server := Server{
		Shortener: shortener,
	}

	return setupRouter(wg, server, authSecret, trustedSubnet)
}
