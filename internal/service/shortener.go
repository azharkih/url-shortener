package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers/models"
	"url-shortener/internal/storage"
)

// Creator определяет поведение сохранения сокращенной ссылки.
type Creator interface {
	CreateShortURL(shortURL *models.ShortURL) (*models.ShortURL, error)
	CreateBatchShortURLs(shortURLs *[]models.ShortURL) error
}

// Retriever определяет поведение извлечения оригинальной ссылки.
type Retriever interface {
	GetShortURL(id string) (*models.ShortURL, error)
}

// Storage объединяет оба интерфейса + проверку соединения с БД
type Storage interface {
	Creator
	Retriever
}

// Service Сервис для работы с URL
type Service struct {
	Repo   Storage
	Config *config.Config
	Logger *zap.SugaredLogger
}

// NewService Конструктор сервиса
func NewService(repo Storage, config *config.Config, logger *zap.SugaredLogger) *Service {
	return &Service{Repo: repo, Config: config, Logger: logger}
}

// PingDB проверяет доступность хранилища
func (s *Service) PingDB(timeoutSeconds ...int) error {
	// Проверяем доступность базы данных
	if dbStorage, ok := s.Repo.(*storage.DatabaseStorage); ok {
		return dbStorage.Ping(timeoutSeconds...)
	}

	// Если это не DatabaseStorage, возвращаем ошибку
	return fmt.Errorf("database storage is not configured")
}

// CreateShortLink Генерация новой короткой ссылки
func (s *Service) CreateShortLink(url string) (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		shortURL := models.NewShortURL(url)

		// Проверка, существует ли уже такая короткая ссылка
		_, err := s.Repo.GetShortURL(shortURL.ID)
		if err != nil {
			res, err := s.Repo.CreateShortURL(shortURL)
			if errors.Is(err, storage.ErrURLAlreadyExists) {
				return fmt.Sprintf("%s/%s", s.Config.BaseShortURL, res.ID), err
			} else if err == nil {
				return fmt.Sprintf("%s/%s", s.Config.BaseShortURL, res.ID), nil
			}
		} else {
			// Если ссылка существует, продолжаем попытки
			continue
		}
	}

	return "", fmt.Errorf("failed to generate short URL after %d attempts", maxAttempts)
}
