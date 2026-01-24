package core

import (
	"context"

	"github.com/karotte128/karotteapi/internal/core"
)

func GetAuthInfo(ctx context.Context) *core.AuthInfo {
	return core.GetAuthInfo(ctx)
}

func HasPermission(info core.AuthInfo, perm string) bool {
	return core.HasPermission(info, perm)
}

func GetModuleConfig(moduleName string) (map[string]any, bool) {
	return core.GetModuleConfig(moduleName)
}

func RegisterMiddleware(middleware core.Middleware) {
	core.RegisterMiddleware(middleware)
}

func RegisterModule(module core.Module) {
	core.RegisterModule(module)
}

func GetNestedValue[Type any](m map[string]any, path ...string) (Type, bool) {
	return core.GetNestedValue[Type](m, path...)
}
