package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// CompressResponse сжимает ответ, если клиент поддерживает GZIP .
func CompressResponse(next http.Handler) http.Handler {
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
		defer func(gzipWriter *gzip.Writer) {
			err := gzipWriter.Close()
			if err != nil {
				log.Println("Error closing gzip writer:", err)
			}
		}(gzipWriter)

		w = &gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (n int, err error) {
	return g.Writer.Write(p)
}
