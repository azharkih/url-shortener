package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"url-shortener/internal/storage"
)

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Invalid request!", http.StatusBadRequest)
		return
	}

	shortURL, err := storage.Repository.ShortURL(id)
	if err != nil {
		http.Error(w, "Not found!", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortURL.FullURL, http.StatusTemporaryRedirect)
}
