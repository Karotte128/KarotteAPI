package core

import (
	"context"
	"slices"

	"github.com/karotte128/karotteapi/apitypes"
)

type authContextKey struct{}

var permissionProvider apitypes.PermissionProvider = nil

func SetAuthInfo(ctx context.Context, info *apitypes.AuthInfo) context.Context {
	return context.WithValue(ctx, authContextKey{}, info)
}

func GetAuthInfo(ctx context.Context) *apitypes.AuthInfo {
	return ctx.Value(authContextKey{}).(*apitypes.AuthInfo)
}

func HasPermission(info apitypes.AuthInfo, perm string) bool {
	contains := slices.Contains(info.Permissions, perm)
	return contains
}

func SetPermissionProvider(provider apitypes.PermissionProvider) {
	permissionProvider = provider
}

func GetPermissionProvider() apitypes.PermissionProvider {
	return permissionProvider
}
