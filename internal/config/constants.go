package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config содержит настройки приложения
type Config struct {
	ServerAddress   string `env:"APP_ADDRESS" envDefault:"localhost:8080"`
	BaseShortURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"shorturls.json"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
}

// NewConfig загружает конфигурацию из .env, переменных окружения и флагов
func NewConfig(sugarLogger *zap.SugaredLogger) (*Config, error) {
	// Загружаем переменные из .env, если файл существует
	if err := godotenv.Load(); err != nil {
		sugarLogger.Warn("Не удалось загрузить .env файл, используем переменные окружения")
	}

	cfg := Config{}

	// Парсим переменные окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	// Проверяем, были ли уже разобраны флаги
	if flag.Parsed() {
		return &cfg, nil
	}

	// Определяем флаги
	serverAddr := flag.String("a", cfg.ServerAddress, "Address for HTTP server (e.g., localhost:8080)")
	baseURL := flag.String("b", cfg.BaseShortURL, "Base URL for short links (e.g., http://localhost:8080)")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "Path to file for storage")
	databaseDSN := flag.String("d", cfg.DatabaseDSN, "DSN for connect to database (e.g., host=localhost port=5432 user=pim password=pim dbname=url_shortener sslmode=disable)")

	// Разбираем флаги
	flag.Parse()

	// Применяем значения флагов, если они были переданы
	if *serverAddr != "" {
		cfg.ServerAddress = *serverAddr
	}
	if *baseURL != "" {
		cfg.BaseShortURL = *baseURL
	}
	if *fileStoragePath != "" {
		cfg.FileStoragePath = *fileStoragePath
	}
	if *databaseDSN != "" {
		cfg.DatabaseDSN = *databaseDSN
	}

	return &cfg, nil
}
