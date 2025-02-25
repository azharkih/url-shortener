package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// UncompressRequest распаковывает запрос, если он сжат.
func UncompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isContentEncoding := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if !isContentEncoding {
			next.ServeHTTP(w, r)
			return
		}
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to read gzip body", http.StatusBadRequest)
			return
		}
		defer func(gzipReader *gzip.Reader) {
			err := gzipReader.Close()
			if err != nil {
				log.Println("Failed to close gzip reader:", err)
			}
		}(gzipReader)

		r.Body = io.NopCloser(gzipReader)
		next.ServeHTTP(w, r)
	})
}
