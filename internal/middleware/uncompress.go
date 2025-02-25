package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// UncompressRequest возвращает middleware, которое распаковывает сжатый запрос.
func UncompressRequest(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Errorw("Failed to read gzip body", "error", err)
				http.Error(w, "Failed to read gzip body", http.StatusBadRequest)
				return
			}
			defer func() {
				if err := gzipReader.Close(); err != nil {
					logger.Errorw("Failed to close gzip reader", "error", err)
				}
			}()

			r.Body = io.NopCloser(gzipReader)
			next.ServeHTTP(w, r)
		})
	}
}
