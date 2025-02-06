package testdata

import (
	"github.com/stretchr/testify/mock"
	"url-shortener/internal/storage"
)

type MockRepository struct {
	mock.Mock
}

// ShortURL возвращает короткую ссылку по ID.
func (m *MockRepository) GetShortURL(id string) (*storage.ShortURL, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.ShortURL), args.Error(1)
}

func (m *MockRepository) SetShortURL(shortURL *storage.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(0)
}

func (m *MockRepository) ShortURLs(ids []string) ([]*storage.ShortURL, error) {
	args := m.Called(ids)
	return args.Get(0).([]*storage.ShortURL), args.Error(1)
}
