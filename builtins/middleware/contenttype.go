package middleware

import (
	"net/http"

	"github.com/karotte128/karotteapi/v2/api"
)

var contentTypeMiddleware = api.Middleware{
	Name:        "contentType",
	Handler:     contentTypeHandler,
	Priority:    2,
	ForceEnable: false,
	Startup:     nil,
	Shutdown:    nil,
}

func contentTypeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	api.RegisterMiddleware(contentTypeMiddleware)
}
