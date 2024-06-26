package util

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("recovered", ErrAttr(err.(error)), slog.Any("trace", debug.Stack()))
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("path", r.URL.String()),
				slog.Int("status", wrapped.status),
				slog.Duration("duration", time.Since(start)))
		}
		return http.HandlerFunc(fn)
	}
}
