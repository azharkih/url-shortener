package app

import (
	"github.com/go-chi/chi/v5"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func NewAppMux() *chi.Mux {
	router := chi.NewRouter()

	// Создаем экземпляр обработчика с зависимостью
	handler := &handlers.Handler{
		CreateShortLink: storage.CreateShortLink, // Передаем функцию генерации ссылок
	}

	// Регистрируем обработчики
	router.Post("/", handler.PostRoot)
	router.Get("/{id}", handlers.GetShortURL)

	return router
}
