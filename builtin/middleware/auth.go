package middleware

import (
	"net/http"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/core"
	"github.com/karotte128/karotteapi/internal"
)

var AuthMiddleware = karotteapi.Middleware{
	Name:     "auth",
	Handler:  AuthHandler,
	Priority: 3,
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("X-API-Key")

		var authInfo karotteapi.AuthInfo

		if header == "" {
			authInfo.ApiKey = ""
			authInfo.Permissions = nil
		} else {
			authInfo.ApiKey = header
			authInfo.Permissions = getPermissions(header)
		}

		ctx := internal.SetAuthInfo(r.Context(), &authInfo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func init() {
	core.RegisterMiddleware(AuthMiddleware)
}

func getPermissions(key string) []string {
	permissionProvider := internal.GetPermissionProvider()
	var permissions []string

	if permissionProvider != nil {
		permissions = permissionProvider(key)
	}

	return permissions
}
