package service

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers/models"
)

// Shortener определяет поведение сохранения сокращенной ссылки.
type Shortener interface {
	SetShortURL(shortURL *models.ShortURL) error
}

// Retriever определяет поведение извлечения оригинальной ссылки.
type Retriever interface {
	GetShortURL(id string) (*models.ShortURL, error)
}

// Storage объединяет оба интерфейса
type Storage interface {
	Shortener
	Retriever
}

// Service Сервис для работы с URL
type Service struct {
	Repo   Storage
	Config *config.Config
	Logger *zap.SugaredLogger
	DB     *sql.DB
}

// NewService Конструктор сервиса
func NewService(repo Storage, config *config.Config, logger *zap.SugaredLogger,
	db *sql.DB) *Service {
	return &Service{Repo: repo, Config: config, Logger: logger, DB: db}
}

// CreateShortLink Генерация новой короткой ссылки
func (s *Service) CreateShortLink(url string) (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		shortURL := models.NewShortURL(url)

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
