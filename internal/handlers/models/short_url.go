package models

import (
	"time"
	"url-shortener/internal/hash"
)

// ShortURL Структура для короткой ссылки
type ShortURL struct {
	ID      string
	FullURL string
	Created int64
}

// NewShortURL Функция для создания объекта ShortURL
func NewShortURL(fullURL string) *ShortURL {
	return &ShortURL{
		ID:      hash.GetRandString(8),
		FullURL: fullURL,
		Created: time.Now().Unix(),
	}
}
