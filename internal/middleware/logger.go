package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestLogger(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w}

			// Логирование запроса
			logger.Infow("Incoming request",
				"method", r.Method,
				"uri", r.RequestURI,
				"time", start,
			)

			// Обработка запроса
			next.ServeHTTP(ww, r)

			// Логирование ответа
			duration := time.Since(start)
			logger.Infow("Response",
				"status", ww.statusCode,
				"response_size", ww.size,
				"duration", duration,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(p)
	rw.size += int64(size)
	return size, err
}
