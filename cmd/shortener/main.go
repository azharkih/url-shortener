package main

import (
	"log"
	"net/http"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

func main() {
	appConfig, err := config.NewConfig()
	if err != nil {
		return
	}
	memoryStorage := storage.NewMemoryStorage()
	storageService := service.NewService(memoryStorage, appConfig)

	log.Printf("Starting server on %s...\n", appConfig.ServerAddress)
	err = http.ListenAndServe(appConfig.ServerAddress, app.NewAppMux(storageService))
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
