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

func (m *MockRepository) SetShortURL(shortURL *models.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(0)
}
