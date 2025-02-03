package storage

import (
	"errors"
	"sync"
)

type MemoryStorage struct {
	sync.RWMutex
	shortURLRecords map[string]*ShortURLRecord
}

type ShortURLRecord struct {
	shortURL *ShortURL
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		sync.RWMutex{},
		map[string]*ShortURLRecord{},
	}
}

func (storage *MemoryStorage) shortURLRecord(idx string) (*ShortURLRecord, error) {
	storage.RLock()
	defer storage.RUnlock()
	if shortURL, ok := storage.shortURLRecords[idx]; ok {
		return shortURL, nil
	}
	return nil, errors.New("short URL not found")
}

func (storage *MemoryStorage) ShortURL(idx string) (*ShortURL, error) {
	if shortURLRecord, err := storage.shortURLRecord(idx); err == nil {
		return shortURLRecord.shortURL, nil
	} else {
		return nil, err
	}
}

func (storage *MemoryStorage) ShortURLs(ids []string) ([]*ShortURL, error) {
	shortURLs := []*ShortURL{}
	for _, idx := range ids {
		if shortURLRecord, err := storage.shortURLRecord(idx); err == nil {
			shortURLs = append(shortURLs, shortURLRecord.shortURL)
		}
	}
	return shortURLs, nil
}

func (storage *MemoryStorage) CreateShortURL(shortURL *ShortURL) error {
	storage.Lock()
	defer storage.Unlock()
	shortURLRecord := ShortURLRecord{shortURL}
	storage.shortURLRecords[shortURL.ID] = &shortURLRecord
	return nil
}

func (storage *MemoryStorage) UpdateShortURL(_ *ShortURL) error {
	return nil
}
