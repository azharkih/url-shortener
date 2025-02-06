package service

import (
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

// Service Сервис для работы с URL
type Service struct {
	Repo   storage.Storage
	Config *config.Config
}

// NewService Конструктор сервиса
func NewService(repo storage.Storage, config *config.Config) *Service {
	return &Service{Repo: repo, Config: config}
}

// CreateShortLink Генерация новой короткой ссылки
func (s *Service) CreateShortLink(url string) (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		shortURL := storage.NewShortURL(url)

		// Проверка, существует ли уже такая короткая ссылка
		_, err := s.Repo.GetShortURL(shortURL.ID)
		if err != nil {
			// Если ошибка, значит ссылки нет
			if err := s.Repo.SetShortURL(shortURL); err == nil {
				return fmt.Sprintf("%s/%s", s.Config.BaseShortURL, shortURL.ID), nil
			}
		} else {
			// Если ссылка существует, продолжаем попытки
			continue
		}
	}

	return "", fmt.Errorf("failed to generate short URL after %d attempts", maxAttempts)
}
