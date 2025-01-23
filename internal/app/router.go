package app

import "net/http"

var router = http.NewServeMux()

func init() {
	router.HandleFunc(`/`, PostRoot)
	router.HandleFunc(`/{id}`, GetShortUrl)
}
