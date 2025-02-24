package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// GzipRequestMiddleware распаковывает запрос, если он сжат.
func UncompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
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
		}
		next.ServeHTTP(w, r)
	})
}
