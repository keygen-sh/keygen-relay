package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/keygen-sh/keygen-relay/internal/logger"
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
		ww := wrapResponseWriter(w)
		start := time.Now()

		next.ServeHTTP(ww, r)

		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"duration", time.Since(start),
		)
	})
}

// signingResponseWriter captures the response body for signing
type signingResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (srw *signingResponseWriter) Write(b []byte) (int, error) {
	return srw.body.Write(b)
}

func (srw *signingResponseWriter) WriteHeader(code int) {
	srw.statusCode = code
}

// SigningMiddleware creates a middleware that signs response bodies with HMAC-SHA256.
// The signature is added as a Relay-Signature header in the format: t=<timestamp>,v1=<signature>
func SigningMiddleware(cfg *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		signer := NewSigner(cfg)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now().Unix()
			w.Header().Set("Relay-Clock", fmt.Sprintf("%d", t))

			if !signer.Enabled() {
				next.ServeHTTP(w, r)

				return
			}

			ww := &signingResponseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(ww, r)

			// signature is computed over "<timestamp>.<raw response body>"
			msg := fmt.Sprintf("%d.%s", t, ww.body.Bytes())
			sig := signer.Sign([]byte(msg))
			w.Header().Set("Relay-Signature", fmt.Sprintf("t=%d,v1=%s", t, hex.EncodeToString(sig)))

			w.WriteHeader(ww.statusCode)
			w.Write(ww.body.Bytes())
		})
	}
}
