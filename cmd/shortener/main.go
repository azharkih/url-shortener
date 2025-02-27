package main

import (
	"log"
	"net/http"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	// Создаем логгер
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

	// Загружаем конфиг
	appConfig, err := config.NewConfig(sugarLogger)
	if err != nil {
		sugarLogger.Fatalw("Failed to load config", "error", err)
	}

	// Подключаемся к базе данных
	db, err := sql.Open("pgx", appConfig.DatabaseDSN)
	if err != nil {
		sugarLogger.Fatalw("Failed to connect to database", "error", err)
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			sugarLogger.Errorf("Failed to close database connection: %v", err)
		}
	}(db)

	// Инициализируем файловое хранилище (если используется)
	fileStorage, err := storage.NewFileStorage(appConfig.FileStoragePath, sugarLogger)
	if err != nil {
		sugarLogger.Fatalw("Failed to initialize file storage", "error", err)
		return
	}

	storageService := service.NewService(fileStorage, appConfig, sugarLogger, db)

	sugarLogger.Infow("Starting server", "address", appConfig.ServerAddress)
	err = http.ListenAndServe(appConfig.ServerAddress, app.NewAppMux(storageService, sugarLogger))
	if err != nil {
		sugarLogger.Fatalw("Server failed", "error", err)
	}
}
