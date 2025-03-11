//go:generate easyjson -all shorten_batch.go

package handlers

import (
	"io"
	"net/http"
	"net/url"
	"url-shortener/internal/handlers/models"
)

//easyjson:json
type ShortenBatchRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

//easyjson:json
type ShortenBatchResponse struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

//easyjson:json
type ShortenBatchRequests []ShortenBatchRequest

//easyjson:json
type ShortenBatchResponses []ShortenBatchResponse

// PostBatchShorten обработчик POST-запроса для создания коротких ссылок по списку
func (h *Handler) PostBatchShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	var requests ShortenBatchRequests
	if err := requests.UnmarshalJSON(body); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var responses ShortenBatchResponses
	var shortURLs []models.ShortURL

	for _, req := range requests {
		if _, err := url.ParseRequestURI(req.URL); err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			h.Service.Logger.Infow("Invalid URL format: %v", err)
			return
		}

		shortURL := models.NewShortURL(req.URL, req.ID)
		shortURLs = append(shortURLs, *shortURL)

		responses = append(responses, ShortenBatchResponse{
			ID:  req.ID,
			URL: h.Service.Config.BaseShortURL + "/" + shortURL.ID,
		})
	}

	if err := h.Service.Repo.CreateBatchShortURLs(&shortURLs); err != nil {
		http.Error(w, "Failed to save batch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseBody, err := responses.MarshalJSON()
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(responseBody)
}
