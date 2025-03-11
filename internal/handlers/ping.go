package handlers

import (
	"net/http"
)

// Ping проверяет доступ к БД и возвращает соответствующий статус
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if err := h.Service.PingDB(1); err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.Service.Logger.Errorf("Error writing response: %v", err)
	}
}
