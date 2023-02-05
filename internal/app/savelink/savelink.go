package savelink

import (
	"github.com/nayakunin/shortener/internal/app/utils"
	"io"
	"net/http"
)

// Storage interface
type Storage interface {
	Get(key string) (string, bool)
	Add(key, link string)
}

// Handler handles POST requests
func Handler(w http.ResponseWriter, r *http.Request, storage Storage) {
	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	key := utils.RandSeq(5)
	url := string(body)

	// add to storage
	storage.Add(key, url)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(key))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
