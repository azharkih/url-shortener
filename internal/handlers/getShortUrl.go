package handlers

import (
	"net/http"
	"url-shortener/internal/storage"
)

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	idx := r.PathValue("id")
	if shortURL, err := storage.Repository.ShortURL(idx); err == nil {
		http.Redirect(w, r, shortURL.FullURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Not found!", http.StatusNotFound)
		return
	}
}
