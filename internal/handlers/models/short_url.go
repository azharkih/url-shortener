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
func NewShortURL(fullURL string, id ...string) *ShortURL {
	// Если id не передан, генерируем новый
	var generatedID string
	if len(id) == 0 || id[0] == "" {
		generatedID = hash.GetRandString(8)
	} else {
		generatedID = id[0]
	}

	return &ShortURL{
		ID:      generatedID,
		FullURL: fullURL,
		Created: time.Now().Unix(),
	}
}
