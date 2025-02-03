package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/storage"
)

func TestGetShortURL(t *testing.T) {
	mockRepo := new(MockRepository)
	storage.Repository = mockRepo // Подмена хранилища

	// Настроим возвращаемое значение для существующей короткой ссылки
	shortURL := &storage.ShortURL{FullURL: "https://example.com", ID: "1234Qwer"}
	mockRepo.On("ShortURL", "1234Qwer").Return(shortURL, nil)
	mockRepo.On("ShortURL", "as123DD1").Return((*storage.ShortURL)(nil), errors.New("not found"))

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		want   want
		path   string
		method string
	}{
		{
			name: "positive case #1",
			want: want{
				code:        307,
				response:    "<a href=\"https://example.com\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
			},
			path:   "/1234Qwer",
			method: http.MethodGet,
		},
		{
			name: "not found",
			want: want{
				code:        404,
				response:    "Not found!\n",
				contentType: "text/plain; charset=utf-8",
			},
			path:   "/as123DD1",
			method: http.MethodGet,
		},
		{
			name: "post method are not allowed",
			want: want{
				code:        405,
				response:    "Only GET requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			path:   "/1234Qwer",
			method: http.MethodPost,
		},
		{
			name: "delete method are not allowed",
			want: want{
				code:        405,
				response:    "Only GET requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			path:   "/1234Qwer",
			method: http.MethodDelete,
		},
		{
			name: "put method are not allowed",
			want: want{
				code:        405,
				response:    "Only GET requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			path:   "/1234Qwer",
			method: http.MethodPut,
		},
		{
			name: "patch method are not allowed",
			want: want{
				code:        405,
				response:    "Only GET requests are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
			path:   "/1234Qwer",
			method: http.MethodPatch,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)

			// Добавляем id в контекст, чтобы обработчик мог его получить
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", test.path[1:]) // Убираем '/'
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()
			GetShortURL(w, request)

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
