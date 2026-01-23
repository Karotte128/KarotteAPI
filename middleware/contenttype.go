package middleware

import (
	"net/http"

	"github.com/karotte128/karotteapi/core"
)

var ContentTypeMiddleware = core.Middleware{
	Name:     "contentType",
	Handler:  ContentTypeHandler,
	Priority: 2,
}

func ContentTypeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	core.RegisterMiddleware(ContentTypeMiddleware)
}
