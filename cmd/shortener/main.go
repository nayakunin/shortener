package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nayakunin/shortener/internal/app/server"
	"github.com/nayakunin/shortener/internal/app/server/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := server.NewRouter()

	log.Fatal(r.Run(fmt.Sprintf("%s:%s", config.Config.ServerAddress, config.Config.Port)))
}
