package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// GetShortURL Обработчик GET-запроса для получения короткой ссылки
func (h *Handler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	// Проверка на правильный метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из URL
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid request!", http.StatusBadRequest)
		return
	}

	// Получаем короткий URL из хранилища
	shortURL, err := h.Service.Repo.GetShortURL(id)
	if err != nil || shortURL == nil {
		// Если не нашли короткую ссылку, возвращаем ошибку
		http.Error(w, "Not found!", http.StatusNotFound)
		return
	}

	// Перенаправление на полный URL
	http.Redirect(w, r, shortURL.FullURL, http.StatusTemporaryRedirect)
}
