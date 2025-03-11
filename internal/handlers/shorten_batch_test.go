package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers/testdata"
	"url-shortener/internal/service"
)

func TestPostBatchShorten(t *testing.T) {
	mockRepo := new(testdata.MockRepository)

	mockRepo.On("CreateBatchShortURLs", mock.Anything).Return(nil).Once()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugarLogger := logger.Sugar()

	cfg, err := config.NewConfig(sugarLogger)
	require.NoError(t, err) // Check that the config was successfully loaded

	mockService := &service.Service{Repo: mockRepo, Config: cfg, Logger: sugarLogger}
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
			body:           `[{"correlation_id": "batch1", "original_url": "https://example1.com"}, {"correlation_id": "batch2", "original_url": "https://example2.com"}]`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `[{"correlation_id":"batch1","short_url":"http://localhost:8080/batch1"},{"correlation_id":"batch2","short_url":"http://localhost:8080/batch2"}]`,
		},
		{
			name:           "negative case #2",
			method:         http.MethodPost,
			body:           "beliberda",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON format\n",
		},
		{
			name:           "negative case #3",
			method:         http.MethodPost,
			body:           `[{"correlation_id": "batch1", "original_url": "example1.com"}, {"correlation_id": "batch2", "original_url": "https://example2.com"}]`,
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
		{
			name:           "empty batch URL",
			method:         http.MethodPost,
			body:           `[{"correlation_id": "batch1", "original_url": ""}]`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL format\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/shorten/batch/", bytes.NewBufferString(test.body))
			w := httptest.NewRecorder()
			handler.PostBatchShorten(w, request)

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
