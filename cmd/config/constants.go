package config

import (
	"flag"
	"fmt"
	"os"
)

// Переменные конфигурации
var (
	ServerAddress string // Адрес HTTP-сервера (-a)
	BaseShortURL  string // Базовый URL сокращённых ссылок (-b)
)

// Инициализация конфигурации
func Init() {
	// Определение флагов
	a := flag.String("a", getEnv("APP_ADDRESS", "localhost:8080"), "Address for HTTP server (e.g., localhost:8888)")
	b := flag.String("b", getEnv("BASE_URL", fmt.Sprintf("http://%s", *a)), "Base URL for short links (e.g., http://localhost:8000/)")

	// Разбор флагов
	flag.Parse()

	// Присваивание значений глобальным переменным
	ServerAddress = *a
	BaseShortURL = *b
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
