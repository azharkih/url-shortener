package handlers

import (
	"github.com/stretchr/testify/mock"
	"url-shortener/internal/storage"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ShortURL(id string) (*storage.ShortURL, error) {
	args := m.Called(id)
	return args.Get(0).(*storage.ShortURL), args.Error(1)
}

func (m *MockRepository) CreateShortURL(shortURL *storage.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(1)
}

func (m *MockRepository) ShortURLs(ids []string) ([]*storage.ShortURL, error) {
	args := m.Called(ids)
	return args.Get(0).([]*storage.ShortURL), args.Error(1)
}

func (m *MockRepository) UpdateShortURL(shortURL *storage.ShortURL) error {
	args := m.Called(shortURL)
	return args.Error(1)
}
