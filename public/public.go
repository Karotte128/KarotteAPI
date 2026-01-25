package public

import (
	"context"

	"github.com/karotte128/karotteapi/apitypes"
	"github.com/karotte128/karotteapi/internal/core"
)

// This function gets the AuthInfo from the request context.
func GetAuthInfo(ctx context.Context) *apitypes.AuthInfo {
	return core.GetAuthInfo(ctx)
}

// This function checks if the AuthInfo of a request has the given permission.
func HasPermission(info apitypes.AuthInfo, perm string) bool {
	return core.HasPermission(info, perm)
}

// This function returns the config of a module.
// It should be used in a module for configurable values.
func GetModuleConfig(moduleName string) (map[string]any, bool) {
	return core.GetModuleConfig(moduleName)
}

// This function should be used inside the init() function of each middleware.
// It adds the middleware to the middleware registry.
func RegisterMiddleware(middleware apitypes.Middleware) {
	core.RegisterMiddleware(middleware)
}

// This function should be used inside the init() function of each module.
// It adds the module to the module registry.
func RegisterModule(module apitypes.Module) {
	core.RegisterModule(module)
}

// This function can be used to get a config value.
// Input the config and the config path.
// Type specifies the type of the return value.
func GetNestedValue[Type any](m map[string]any, path ...string) (Type, bool) {
	return core.GetNestedValue[Type](m, path...)
}
