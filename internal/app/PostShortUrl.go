package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

func PostRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fullUrl := string(data)
	_, err = url.ParseRequestURI(fullUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortUrl := createShortUrl(fullUrl)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(shortUrl))
}

func createShortUrl(url string) string {
	shortcode := getRandString(8)
	for {
		if _, ok := Urls[shortcode]; !ok {
			Urls[shortcode] = url
			return fmt.Sprintf("http://%s/%s", BaseUrl, shortcode)
		}
	}
}

func getRandString(countOfChars int) string {
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var shortcode strings.Builder

	for i := 0; i < countOfChars; i++ {
		random := rand.Intn(len(charSet))
		shortcode.WriteString(string(charSet[random]))
	}

	return shortcode.String()
}
