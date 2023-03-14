package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"
	storagePackage "github.com/nayakunin/shortener/internal/app/storage"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storagePackage.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(*cfg, storage)
	r := server.NewRouter(s)

	log.Fatal(r.Run(cfg.ServerAddress))
}
