package public

import (
	"context"

	"github.com/karotte128/karotteapi/apitypes"
	"github.com/karotte128/karotteapi/internal/core"
)

func GetAuthInfo(ctx context.Context) *apitypes.AuthInfo {
	return core.GetAuthInfo(ctx)
}

func HasPermission(info apitypes.AuthInfo, perm string) bool {
	return core.HasPermission(info, perm)
}

func GetModuleConfig(moduleName string) (map[string]any, bool) {
	return core.GetModuleConfig(moduleName)
}

func RegisterMiddleware(middleware apitypes.Middleware) {
	core.RegisterMiddleware(middleware)
}

func RegisterModule(module apitypes.Module) {
	core.RegisterModule(module)
}

func GetNestedValue[Type any](m map[string]any, path ...string) (Type, bool) {
	return core.GetNestedValue[Type](m, path...)
}
