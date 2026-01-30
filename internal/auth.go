package internal

import (
	"context"
	"log"
	"slices"

	"github.com/karotte128/karotteapi"
)

type authContextKey struct{}

var permissionProvider karotteapi.PermissionProvider = nil

func SetAuthInfo(ctx context.Context, info *karotteapi.AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey{}, info)
}

func GetAuthInfo(ctx context.Context) *karotteapi.AuthInfo {
	return ctx.Value(authContextKey{}).(*karotteapi.AuthInfo)
}

func HasPermission(info karotteapi.AuthInfo, perm string) bool {
	contains := slices.Contains(info.Permissions, perm)
	return contains
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
