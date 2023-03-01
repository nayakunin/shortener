package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"
	"github.com/nayakunin/shortener/internal/app/storage"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	a := os.Getenv("FILE_STORAGE_PATH")

	fmt.Println(a)

	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New()
	defer func(s *storage.Storage) {
		err := s.Close()
		log.Fatal(err)
	}(s)

	r := server.NewRouter(s)

	log.Fatal(r.Run(config.Config.ServerAddress))
}
