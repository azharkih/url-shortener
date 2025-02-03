package app

import (
	"net/http"
	"url-shortener/internal"
)

func Run() error {
	router := NewAppMux()
	return http.ListenAndServe(internal.BaseURL, router)
}
