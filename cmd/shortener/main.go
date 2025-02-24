package main

import (
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
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	appConfig, err := config.NewConfig()
	if err != nil {
		sugarLogger.Fatalw("Failed to load config", "error", err)
		return
	}

	memoryStorage := storage.NewMemoryStorage()
	storageService := service.NewService(memoryStorage, appConfig)

	sugarLogger.Infow("Starting server", "address", appConfig.ServerAddress)

	err = http.ListenAndServe(appConfig.ServerAddress, app.NewAppMux(storageService, sugarLogger))
	if err != nil {
		sugarLogger.Fatalw("Server failed", "error", err)
	}
}
