package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// CompressResponse возвращает middleware, которое сжимает ответ, если клиент поддерживает GZIP.
func CompressResponse(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			isAcceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
			isAvailableType := strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")

			if !isAcceptEncoding || !isAvailableType {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer func() {
				if err := gzipWriter.Close(); err != nil {
					logger.Errorw("Error closing gzip writer", "error", err)
				}
			}()

			w = &gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}
			next.ServeHTTP(w, r)
		})
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.Writer.Write(p)
}
