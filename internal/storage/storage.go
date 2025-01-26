package storage

import (
	"math/rand"
	"strings"
	"time"
)

var Repository = NewMemoryStorage()

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
