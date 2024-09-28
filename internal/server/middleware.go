package server

import (
	"log/slog"
	"net/http"
	"time"
)

func wrapResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) Status() int {
	return lrw.statusCode
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedResponse := wrapResponseWriter(w)

		next.ServeHTTP(wrappedResponse, r)

		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedResponse.Status(),
			"remote_addr", r.RemoteAddr,
			"duration", time.Since(start),
		)
	})
}
