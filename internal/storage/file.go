package storage

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"url-shortener/internal/handlers/models"
)

// FileStorage реализует хранение в файле
type FileStorage struct {
	filePath string
	data     map[string]*models.ShortURL
	logger   *zap.SugaredLogger
}

// NewFileStorage создаёт новое хранилище для файла
func NewFileStorage(filePath string, logger *zap.SugaredLogger) (*FileStorage, error) {
	fs := &FileStorage{
		filePath: filePath,
		data:     make(map[string]*models.ShortURL),
		logger:   logger,
	}

	// Попробуем загрузить данные из файла, если они есть
	if err := fs.load(); err != nil {
		return nil, fmt.Errorf("failed to load file storage: %w", err)
	}

	return fs, nil
}

// load загружает данные из файла
func (fs *FileStorage) load() error {
	file, err := os.Open(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, просто возвращаем пустое хранилище
			return nil
		}
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fs.logger.Errorf("Failed to close file: %v", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var shortURL models.ShortURL
		if err := decoder.Decode(&shortURL); err != nil {
			return err
		}
		fs.data[shortURL.ID] = &shortURL
	}

	return nil
}

// save сохраняет данные в файл
func (fs *FileStorage) save() error {
	file, err := os.Create(fs.filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fs.logger.Errorf("Failed to close file: %v", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	for _, shortURL := range fs.data {
		if err := encoder.Encode(shortURL); err != nil {
			return err
		}
	}

	return nil
}

// SetShortURL сохраняет сокращённый URL в файл
func (fs *FileStorage) SetShortURL(shortURL *models.ShortURL) error {
	fs.data[shortURL.ID] = shortURL
	return fs.save()
}

// GetShortURL возвращает сокращённый URL по ID
func (fs *FileStorage) GetShortURL(id string) (*models.ShortURL, error) {
	shortURL, exists := fs.data[id]
	if !exists {
		return nil, fmt.Errorf("short URL not found")
	}
	return shortURL, nil
}
