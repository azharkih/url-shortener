package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v10"
)

// Config содержит настройки приложения
type Config struct {
	ServerAddress string `env:"APP_ADDRESS" envDefault:"localhost:8080"`
	BaseShortURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

// NewConfig загружает конфигурацию из переменных окружения и флагов
func NewConfig() (*Config, error) {
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
	serverAddr := flag.String(
		"a", cfg.ServerAddress, "Address for HTTP server (e.g., localhost:8080)")
	baseURL := flag.String(
		"b", cfg.BaseShortURL, "Base URL for short links (e.g., http://localhost:8080)")

	// Разбираем флаги
	flag.Parse()

	// Применяем значения флагов, если они были переданы
	if *serverAddr != "" {
		cfg.ServerAddress = *serverAddr
	}
	if *baseURL != "" {
		cfg.BaseShortURL = *baseURL
	} else {
		cfg.BaseShortURL = fmt.Sprintf("http://%s", cfg.ServerAddress)
	}

	return &cfg, nil
}
