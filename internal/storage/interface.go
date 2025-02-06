package storage

import (
	"time"
	"url-shortener/internal/utils"
)

// Storage  Интерфейс хранилища
type Storage interface {
	GetShortURL(id string) (*ShortURL, error)
	SetShortURL(shortURL *ShortURL) error
}

// ShortURL Структура для короткой ссылки
type ShortURL struct {
	ID      string
	FullURL string
	Created int64
}

// NewShortURL Функция для создания объекта ShortURL
func NewShortURL(fullURL string) *ShortURL {
	return &ShortURL{
		ID:      utils.GetRandString(8),
		FullURL: fullURL,
		Created: time.Now().Unix(),
	}
}
