package storage

import "url-shortener/internal/handlers/models"

// TODO переместить в другой пакет

// Shortener определяет поведение сохранения сокращенной ссылки.
type Shortener interface {
	SetShortURL(shortURL *models.ShortURL) error
}

// Retriever определяет поведение извлечения оригинальной ссылки.
type Retriever interface {
	GetShortURL(id string) (*models.ShortURL, error)
}
