package handlers

import "url-shortener/internal/storage"

type Handler struct {
	Service *storage.Service
}
