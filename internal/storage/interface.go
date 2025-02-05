package storage

import (
	"fmt"
	"math/rand"
	"time"
	"url-shortener/internal/config"
)

// Storage  Интерфейс хранилища
type Storage interface {
	ShortURL(id string) (*ShortURL, error)
	ShortURLs(ids []string) ([]*ShortURL, error)
	CreateShortURL(shortURL *ShortURL) error
	UpdateShortURL(shortURL *ShortURL) error
}

// ShortURL Структура для короткой ссылки
type ShortURL struct {
	ID      string
	FullURL string
	Created int64
}

// Service Сервис для работы с URL
type Service struct {
	Repo Storage
}

// NewService Конструктор сервиса
func NewService(repo Storage) *Service {
	return &Service{Repo: repo}
}

// CreateShortLink Генерация новой короткой ссылки
func (s *Service) CreateShortLink(url string) (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		shortURL := NewShortURL(url)

		// Проверка, существует ли уже такая короткая ссылка
		_, err := s.Repo.ShortURL(shortURL.ID)
		if err != nil {
			// Если ошибка, значит ссылки нет
			if err := s.Repo.CreateShortURL(shortURL); err == nil {
				return fmt.Sprintf("%s/%s", config.BaseShortURL, shortURL.ID), nil
			}
		} else {
			// Если ссылка существует, продолжаем попытки
			continue
		}
	}

	return "", fmt.Errorf("failed to generate short URL after %d attempts", maxAttempts)
}

// NewShortURL Функция для создания объекта ShortURL
func NewShortURL(fullURL string) *ShortURL {
	return &ShortURL{
		ID:      getRandString(8),
		FullURL: fullURL,
		Created: time.Now().Unix(),
	}
}

// Генерация случайной строки
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandString(countOfChars int) string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, countOfChars)
	for i := range bytes {
		bytes[i] = charSet[rng.Intn(len(charSet))]
	}
	return string(bytes)
}
