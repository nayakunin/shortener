package getlink

import (
	"net/http"
)

// Storage interface
type Storage interface {
	Get(key string) (string, bool)
	Add(key, link string)
}

// Handler handles GET requests
func Handler(w http.ResponseWriter, r *http.Request, storage Storage) {
	id := r.URL.Path[1:]
	if id == "" {
		// throw error
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	link, ok := storage.Get(id)
	if !ok {
		// throw error
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}
