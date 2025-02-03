package storage

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"url-shortener/internal"
)

var Repository Storage = NewMemoryStorage()

type Storage interface {
	ShortURL(id string) (*ShortURL, error)
	ShortURLs(ids []string) ([]*ShortURL, error)
	CreateShortURL(shortURL *ShortURL) error
	UpdateShortURL(shortURL *ShortURL) error
}

type ShortURL struct {
	ID      string
	FullURL string
	Created int64
}

func NewShortURL(FullURL string) *ShortURL {
	now := time.Now().Unix()
	shortURL := ShortURL{
		ID:      getRandString(8),
		FullURL: FullURL,
		Created: now,
	}
	return &shortURL
}

func getRandString(countOfChars int) string {
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var shortcode strings.Builder

	for i := 0; i < countOfChars; i++ {
		random := rand.Intn(len(charSet))
		shortcode.WriteString(string(charSet[random]))
	}

	return shortcode.String()
}

// Функция генерации ссылки
func CreateShortLink(url string) string {
	for {
		shortURL := NewShortURL(url)
		if item, _ := Repository.ShortURL(shortURL.ID); item != nil {
			continue
		}
		if err := Repository.CreateShortURL(shortURL); err == nil {
			return fmt.Sprintf("http://%s/%s", internal.BaseURL, shortURL.ID)
		}
	}
}
