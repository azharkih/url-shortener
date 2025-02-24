package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

//easyjson:json
type ShortenRequest struct {
	URL string `json:"url"`
}

//easyjson:json
type ShortenResponse struct {
	Result string `json:"result"`
}

// PostShorten обработчик POST-запроса для создания короткой ссылки через JSON
func (h *Handler) PostShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Error reading request body: %v", err)
		return
	}
	defer r.Body.Close()

	var req ShortenRequest
	if err := req.UnmarshalJSON(data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		log.Printf("Invalid JSON: %v", err)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		log.Printf("Invalid URL format: %v", err)
		return
	}

	shortURL, err := h.Service.CreateShortLink(req.URL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		log.Printf("Error creating short URL: %v", err)
		return
	}

	resp := ShortenResponse{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonData, err := resp.MarshalJSON()
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		log.Printf("Error marshaling response: %v", err)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
