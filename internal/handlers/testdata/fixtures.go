package testdata

import (
	"github.com/stretchr/testify/mock"
	"url-shortener/internal/handlers/models"
)

type MockRepository struct {
	mock.Mock
}

// GetShortURL возвращает короткую ссылку по ID.
func (m *MockRepository) GetShortURL(id string) (*models.ShortURL, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShortURL), args.Error(1)
}

// CreateShortURL возвращает ошибку, если URL уже существует.
func (m *MockRepository) CreateShortURL(shortURL *models.ShortURL) (*models.ShortURL, error) {
	args := m.Called(shortURL)
	return shortURL, args.Error(0)
}

func (m *MockRepository) CreateBatchShortURLs(shortURLs *[]models.ShortURL) error {
	args := m.Called(shortURLs)
	return args.Error(0)
}
