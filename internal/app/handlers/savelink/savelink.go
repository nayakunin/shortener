package savelink

import (
	"fmt"
	"io"
	"net/http"
)

// Storage interface
type Storage interface {
	Add(link string) string
}

// Handler handles POST requests
func Handler(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		url := string(body)

		// add to storage
		key := s.Add(url)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", key)))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

}
