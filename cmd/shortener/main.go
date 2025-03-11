package main

import (
	"database/sql"
	"net/http"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	// Создаем логгер
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	sugarLogger := logger.Sugar()

	// Загружаем конфиг
	appConfig, err := config.NewConfig(sugarLogger)
	if err != nil {
		sugarLogger.Fatalw("Failed to load config", "error", err)
	}

	var appStorage service.Storage

	// Проверяем доступность БД
	if appConfig.DatabaseDSN != "" {
		db, err := sql.Open("pgx", appConfig.DatabaseDSN)
		if err != nil {
			sugarLogger.Warnw("Failed to connect to database, falling back to file storage", "error", err)
		} else {
			sugarLogger.Infow("Connected to database")
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					sugarLogger.Errorf("Failed to close database connection: %v", err)
				}
			}(db)
			appStorage = storage.NewDatabaseStorage(db, sugarLogger)
		}
	}

	// Если БД не доступна, проверяем файловое хранилище
	if appStorage == nil && appConfig.FileStoragePath != "" {
		fileStorage, err := storage.NewFileStorage(appConfig.FileStoragePath, sugarLogger)
		if err != nil {
			sugarLogger.Warnw("Failed to initialize file storage, falling back to in-memory storage", "error", err)
		} else {
			sugarLogger.Infow("Using file storage")
			appStorage = fileStorage
		}
	}

	// Если файловое хранилище тоже не доступно, используем память
	if appStorage == nil {
		sugarLogger.Infow("Using in-memory storage")
		appStorage = storage.NewMemoryStorage()
	}

	appService := service.NewService(appStorage, appConfig, sugarLogger)

	// Запускаем сервер
	sugarLogger.Infow("Starting server", "address", appConfig.ServerAddress)
	err = http.ListenAndServe(appConfig.ServerAddress, app.NewAppMux(appService, sugarLogger))
	if err != nil {
		sugarLogger.Fatalw("Server failed", "error", err)
	}
}
