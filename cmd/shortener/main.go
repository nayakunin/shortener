package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New()

	r := server.NewRouter(s)

	log.Fatal(r.Run(config.Config.ServerAddress))
}
