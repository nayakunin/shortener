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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	s, err := storage.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := server.NewRouter(*cfg, s)

	log.Fatal(r.Run(cfg.ServerAddress))
}
