package main

import (
	"log"
	"net/http"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}(logger)
	sugarLogger := logger.Sugar()

	appConfig, err := config.NewConfig()
	if err != nil {
		sugarLogger.Fatalw("Failed to load config", "error", err)
	}

	fileStorage, err := storage.NewFileStorage(appConfig.FileStoragePath, sugarLogger)
	if err != nil {
		sugarLogger.Fatalw("Failed to initialize file storage", "error", err)
		return
	}

	storageService := service.NewService(fileStorage, appConfig, sugarLogger)

	sugarLogger.Infow("Starting server", "address", appConfig.ServerAddress)
	err = http.ListenAndServe(appConfig.ServerAddress, app.NewAppMux(storageService, sugarLogger))
	if err != nil {
		sugarLogger.Fatalw("Server failed", "error", err)
	}
}
