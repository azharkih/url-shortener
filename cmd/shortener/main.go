package main

import (
	"log"
	"net/http"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

func main() {
	config.Init()
	memoryStorage := storage.NewMemoryStorage()
	storageService := storage.NewService(memoryStorage)

	log.Printf("Starting server on %s...\n", config.ServerAddress)
	err := http.ListenAndServe(config.ServerAddress, app.NewAppMux(storageService))
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
