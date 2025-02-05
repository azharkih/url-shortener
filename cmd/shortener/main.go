package main

import (
	"log"
	"net/http"
	"url-shortener/cmd/config"
	"url-shortener/internal/app"
)

func main() {
	// Инициализация конфигурации
	config.Init()

	log.Printf("Starting server on %s...\n", config.ServerAddress)
	err := http.ListenAndServe(config.ServerAddress, app.NewAppMux())
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
