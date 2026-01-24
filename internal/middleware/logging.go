package middleware

import (
	"log"
	"net/http"

	"github.com/karotte128/karotteapi/apitypes"
	"github.com/karotte128/karotteapi/internal/core"
)

var LoggingMiddleware = apitypes.Middleware{
	Name:     "logging",
	Handler:  LoggingHandler,
	Priority: 1,
}

// loggingResponseWriter wraps http.ResponseWriter so we can capture
// the status code and response size written by handlers.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// LoggingMiddleware logs method, path, status, and response size.
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default, in case WriteHeader is never called
		}

		next.ServeHTTP(lrw, r)

		log.Printf("[LOG] %s %s -> %d (%d bytes)",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			lrw.size,
		)
	})
}

// Automatically register this middleware globally
func init() {
	core.RegisterMiddleware(LoggingMiddleware)
}
