package app

import (
	"github.com/go-chi/chi/v5"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func NewAppMux(s *storage.Service) *chi.Mux {
	router := chi.NewRouter()

	handler := &handlers.Handler{
		Service: s,
	}

	// Регистрируем обработчики
	router.Post("/", handler.PostRoot)
	router.Get("/{id}", handler.GetShortURL)

	return router
}
