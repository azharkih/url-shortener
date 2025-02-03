package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/storage"
)

func TestPostRoot(t *testing.T) {
	mockRepo := new(MockRepository)
	storage.Repository = mockRepo // Подмена хранилища

	// Настроим возвращаемое значение для существующей короткой ссылки
	shortURL := &storage.ShortURL{FullURL: "https://example.com", ID: "mock1234"}
	mockRepo.On("ShortURL", "mock1234").Return(shortURL, nil)
	mockRepo.On("CreateShortURL", mock.Anything).Return(nil)

	mockCreateShort := func(url string) string {
		return "http://short.url/mock1234"
	}

	handler := &Handler{CreateShortLink: mockCreateShort}

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		want   want
		method string
	}{
		{
			name: "positive case #1",
			want: want{
				code:        201,
				response:    "http://short.url/mock1234",
				contentType: "text/plain",
			},
			method: http.MethodPost,
		},
		{
			name: "post method are not allowed",
			want: want{
				code:        405,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			method: http.MethodGet,
		},
		{
			name: "delete method are not allowed",
			want: want{
				code:        405,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			method: http.MethodDelete,
		},
		{
			name: "put method are not allowed",
			want: want{
				code:        405,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			method: http.MethodPut,
		},
		{
			name: "patch method are not allowed",
			want: want{
				code:        405,
				response:    "Only POST requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			method: http.MethodPatch,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fullURL := "https://example.com"
			request := httptest.NewRequest(test.method, "/", bytes.NewBufferString(fullURL))

			w := httptest.NewRecorder()
			handler.PostRoot(w, request)

			res := w.Result()
			err := res.Body.Close()
			require.NoError(t, err)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
