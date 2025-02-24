package app

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"url-shortener/internal/handlers"
	"url-shortener/internal/middleware"
	"url-shortener/internal/service"
)

func NewAppMux(s *service.Service, logger *zap.SugaredLogger) *chi.Mux {
	router := chi.NewRouter()
	handler := handlers.NewHandler(s)

	// Регистрируем middleware для логирования
	router.Use(middleware.RequestLogger(logger))

	// Регистрируем обработчики
	router.Post("/", handler.PostRoot)
	router.Get("/{id}", handler.GetShortURL)

	return router
}
