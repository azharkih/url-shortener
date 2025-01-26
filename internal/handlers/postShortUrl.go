package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"url-shortener/internal"
	"url-shortener/internal/storage"
)

func PostRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
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
	link := createShortLink(fullURL)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(link))
}

func createShortLink(url string) string {
	for {
		shortURL := storage.NewShortURL(url)
		if item, _ := storage.Repository.ShortURL(shortURL.ID); item != nil {
			continue
		}
		if err := storage.Repository.CreateShortURL(shortURL); err == nil {
			return fmt.Sprintf("http://%s/%s", internal.BaseURL, shortURL.ID)
		}

	}
}
