package app

import (
	"net/http"
	"url-shortener/internal/handlers"
)

var router = http.NewServeMux()

func init() {
	router.HandleFunc(`/`, handlers.PostRoot)
	router.HandleFunc(`/{id}`, handlers.GetShortURL)
}
