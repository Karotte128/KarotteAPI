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
