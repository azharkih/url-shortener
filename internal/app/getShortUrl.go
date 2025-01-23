package app

import "net/http"

func GetShortUrl(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	idx := r.PathValue("id")
	fullUrl, ok := Urls[idx]
	if !ok {
		http.Error(w, "Not found!", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fullUrl, http.StatusTemporaryRedirect)
}
