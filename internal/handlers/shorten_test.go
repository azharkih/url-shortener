package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers/models"
	"url-shortener/internal/handlers/testdata"
	"url-shortener/internal/service"
)

func TestShorten(t *testing.T) {
	mockRepo := new(testdata.MockRepository)

	mockRepo.On("GetShortURL", mock.Anything).Return(nil, errors.New("not found")).Maybe()

	mockRepo.On("SetShortURL", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		shortURL := args.Get(0).(*models.ShortURL)
		shortURL.ID = "mock1234"
	}).Once()

	cfg, err := config.NewConfig()
	require.NoError(t, err) // Проверяем, что конфиг успешно загружен

	mockService := &service.Service{Repo: mockRepo, Config: cfg}
	handler := NewHandler(mockService)

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "positive case #1",
			method:         http.MethodPost,
			body:           "https://example.com",
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf("%s/%s", mockService.Config.BaseShortURL, "mock1234"),
		},
		{
			name:           "negative case #2",
			method:         http.MethodPost,
			body:           "beliberda",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL format\n",
		},
		{
			name:           "only POST method allowed",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Only POST requests are allowed!\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", bytes.NewBufferString(test.body))
			w := httptest.NewRecorder()
			handler.PostRoot(w, request)

			res := w.Result()
			err := res.Body.Close()
			require.NoError(t, err)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, test.expectedBody, string(resBody))
		})
	}

	// Проверяем, что все ожидания моков выполнены
	mockRepo.AssertExpectations(t)
}
