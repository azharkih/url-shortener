package handlers

import (
	"io"
	"net/http"
	"net/url"
	"url-shortener/internal/storage"
)

// Handler содержит зависимости
type Handler struct {
	CreateShortLink func(string) string
	Repo            storage.Storage
}

func (h *Handler) PostRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fullURL := string(data)
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link := h.CreateShortLink(fullURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(link))
}
