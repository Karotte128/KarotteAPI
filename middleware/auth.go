package middleware

import (
	"net/http"

	"github.com/karotte128/karotteapi/core"
)

var AuthMiddleware = core.Middleware{
	Name:     "auth",
	Handler:  AuthHandler,
	Priority: 3,
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("X-API-Key")

		var authInfo core.AuthInfo

		if header == "" {
			authInfo.ApiKey = ""
			authInfo.Permissions = nil
		} else {
			authInfo.ApiKey = header
			authInfo.Permissions = getPermissions(header)
		}

		ctx := core.SetAuthInfo(r.Context(), &authInfo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func init() {
	core.RegisterMiddleware(AuthMiddleware)
}

func getPermissions(key string) []string {
	var permissions []string

	//TODO: add permissions system

	return permissions
}
