package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nayakunin/shortener/internal/app/server"
)

const PORT = ":8080"

func main() {
	rand.Seed(time.Now().UnixNano())

	r := server.NewRouter()

	log.Fatal(r.Run(PORT))
}
