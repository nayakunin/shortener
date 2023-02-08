package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nayakunin/shortener/internal/app/handlers/getlink"
	"github.com/nayakunin/shortener/internal/app/handlers/savelink"
	"github.com/nayakunin/shortener/internal/app/storage"
	"log"
	"math/rand"
	"time"
)

func setupRouter(s *storage.Storage) *gin.Engine {
	r := gin.Default()
	r.POST("/", savelink.Handler(s))
	r.GET("/:id", getlink.Handler(s))
	return r
}

func main() {
	rand.Seed(time.Now().UnixNano())
	s := storage.New()

	r := setupRouter(s)

	log.Fatal(r.Run(":8080"))
}
