package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// PostRoot Обработчик POST-запроса для создания короткой ссылки
func (h *Handler) PostRoot(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Error reading request body: %v", err)
		return
	}
	fullURL := string(data)

	// Проверка правильности URL
	_, err = url.ParseRequestURI(fullURL)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		log.Printf("Invalid URL format: %v", err)
		return
	}

	// Генерация короткой ссылки через сервис
	link, err := h.Service.CreateShortLink(fullURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		log.Printf("Error creating short URL: %v", err)
		return
	}

	// Отправка результата
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(link))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
