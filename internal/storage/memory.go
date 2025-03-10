package storage

import (
	"errors"
	"sync"
	"url-shortener/internal/handlers/models"
)

// MemoryStorage Структура для хранения в памяти
type MemoryStorage struct {
	sync.RWMutex
	shortURLRecords map[string]*ShortURLRecord
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortURLRecords: make(map[string]*ShortURLRecord),
	}
}

type ShortURLRecord struct {
	shortURL *models.ShortURL
}

func (storage *MemoryStorage) shortURLRecord(idx string) (*ShortURLRecord, error) {
	storage.RLock()
	defer storage.RUnlock()
	if shortURL, ok := storage.shortURLRecords[idx]; ok {
		return shortURL, nil
	}
	return nil, errors.New("short URL not found")
}

func (storage *MemoryStorage) GetShortURL(idx string) (*models.ShortURL, error) {
	if shortURLRecord, err := storage.shortURLRecord(idx); err == nil {
		return shortURLRecord.shortURL, nil
	} else {
		return nil, err
	}
}

func (storage *MemoryStorage) SetShortURL(shortURL *models.ShortURL) error {
	storage.Lock()
	defer storage.Unlock()
	shortURLRecord := ShortURLRecord{shortURL}
	storage.shortURLRecords[shortURL.ID] = &shortURLRecord
	return nil
}
