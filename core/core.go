package core

import (
	"context"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/internal"
)

// This function returns the config of a module.
// It should be used in a module for configurable values.
func GetModuleConfig(moduleName string) (map[string]any, bool) {
	return internal.GetModuleConfig(moduleName)
}

// This function should be used inside the init() function of each middleware.
// It adds the middleware to the middleware registry.
func RegisterMiddleware(middleware karotteapi.Middleware) {
	internal.RegisterMiddleware(middleware)
}

// This function should be used inside the init() function of each module.
// It adds the module to the module registry.
func RegisterModule(module karotteapi.Module) {
	internal.RegisterModule(module)
}

// This function can be used to get a config value.
// Input the config and the config path.
// Type specifies the type of the return value.
func GetNestedValue[Type any](m map[string]any, path ...string) (Type, bool) {
	return internal.GetNestedValue[Type](m, path...)
}

// This function adds additional info to the request context.
// It is usually used by a middleware.
func SetRequestContext(ctx context.Context, info *karotteapi.RequestContext) context.Context {
	return internal.SetRequestContext(ctx, info)
}

// This function retrieves the additional info from the request context.
// It is usually used in a module.
func GetRequestContext(ctx context.Context, contextKey string) karotteapi.RequestContext {
	return internal.GetRequestContext(ctx, contextKey)
}
