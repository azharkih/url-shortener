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

	// Регистрируем middleware
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.CompressResponse(logger))
	router.Use(middleware.UncompressRequest(logger))

	// Регистрируем обработчики
	router.Post("/", handler.PostRoot)
	router.Get("/{id}", handler.GetShortURL)
	router.Post("/api/shorten", handler.PostShorten)

	return router
}
