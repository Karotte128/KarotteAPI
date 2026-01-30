package internal

import (
	"context"
	"log"

	"github.com/karotte128/karotteapi"
)

// AuthInfo is created by the auth middleware.
// It contains the authentication status and permissions of the request.
type AuthInfo struct {
	// ApiKey is the raw key sent by the user. Do not use this.
	ApiKey string

	// Permissions is the list of permissions the user has.
	Permissions []string
}

type authContextKey struct{}

var permissionProvider karotteapi.PermissionProvider = nil

func SetAuthInfo(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey{}, info)
}

func GetAuthInfo(ctx context.Context) *AuthInfo {
	return ctx.Value(authContextKey{}).(*AuthInfo)
}

func SetPermissionProvider(provider karotteapi.PermissionProvider) {
	if provider == nil {
		log.Println("[AUTHENTICATION] No permission provider registered!")
	}

	permissionProvider = provider
}

func GetPermissionProvider() karotteapi.PermissionProvider {
	return permissionProvider
}
