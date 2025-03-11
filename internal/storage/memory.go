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

func (storage *MemoryStorage) CreateShortURL(shortURL *models.ShortURL) error {
	storage.Lock()
	defer storage.Unlock()
	shortURLRecord := ShortURLRecord{shortURL}
	storage.shortURLRecords[shortURL.ID] = &shortURLRecord
	return nil
}

// CreateBatchShortURLs сохраняет список сокращенных URL в памяти
func (storage *MemoryStorage) CreateBatchShortURLs(shortURLs *[]models.ShortURL) error {
	if len(*shortURLs) == 0 {
		return nil // Нет данных для сохранения
	}

	storage.Lock()
	defer storage.Unlock()

	for i := range *shortURLs {
		storage.shortURLRecords[(*shortURLs)[i].ID] = &ShortURLRecord{shortURL: &(*shortURLs)[i]}
	}

	return nil
}
