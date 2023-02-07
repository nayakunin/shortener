package getlink

import (
	"net/http"
)

// Storage interface
type Storage interface {
	Get(key string) (string, bool)
}

// Handler handles GET requests
func Handler(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[1:]
		if id == "" {
			// throw error
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		link, ok := s.Get(id)
		if !ok {
			// throw error
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	}
}
