package core

import (
	"net/http"

	cfg "github.com/karotte128/karottelib/config"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/internal"
)

// This function returns the config of a module.
// It should be used in a module for configurable values.
func GetModuleConfig(moduleName string) (karotteapi.Config, bool) {
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
func GetNestedValue[Type any](m karotteapi.Config, path ...string) (Type, bool) {
	return cfg.GetNestedValue[Type](m, path...)
}

// This function adds additional info to the request context.
// It is usually used by a middleware.
func SetRequestContext(r *http.Request, key string, value any) {
	internal.SetRequestContext(r, key, value)
}

// This function retrieves the additional info from the request context.
// It is usually used in a module.
func GetRequestContext[T any](r *http.Request, key string) (value T, ok bool) {
	return internal.GetRequestContext[T](r, key)
}
