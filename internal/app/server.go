package app

import (
	"net/http"
	"url-shortener/internal"
)

func Run() error {
	return http.ListenAndServe(internal.BaseURL, router)
}
