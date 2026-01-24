package core

import (
	"context"
	"slices"
)

type authContextKey struct{}

type AuthInfo struct {
	ApiKey      string
	Permissions []string
}

type PermissionProvider func(key string) []string

var permissionProvider PermissionProvider = nil

func SetAuthInfo(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey{}, info)
}

func GetAuthInfo(ctx context.Context) *AuthInfo {
	return ctx.Value(authContextKey{}).(*AuthInfo)
}

func HasPermission(info AuthInfo, perm string) bool {
	contains := slices.Contains(info.Permissions, perm)
	return contains
}

func SetPermissionProvider(provider PermissionProvider) {
	permissionProvider = provider
}

func GetPermissionProvider() PermissionProvider {
	return permissionProvider
}
