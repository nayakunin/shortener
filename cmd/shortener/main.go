package main

import (
	"github.com/nayakunin/shortener/internal/app/handlers/getlink"
	"github.com/nayakunin/shortener/internal/app/handlers/savelink"
	"github.com/nayakunin/shortener/internal/app/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s := storage.New()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			savelink.Handler(s)(w, r)
		case http.MethodGet:
			getlink.Handler(s)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
