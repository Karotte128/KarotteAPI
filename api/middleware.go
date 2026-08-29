package api

import (
	"log"
	"net/http"
	"sort"

	cfg "github.com/karotte128/karottelib/config"
)

// Middleware is a function that wraps an http.Handler and returns a new one.
// This allows transforming the request/response pipeline.
//
// Examples:
// - logging
// - authentication
// - rate limiting

// Middlewares register themselves automatically via init() inside their
// own package. The core does not need to know about them explicitly.

// Middleware is the struct the middleware needs to provide to the middleware registry to register itself.
type Middleware struct {
	// Name is the name of the middleware. It is used for logging.
	Name string

	// Priority is the order in which middlewares should be registered.
	// Lower number means the middleware gets registered earlier (higher priority).
	Priority uint

	// Middleware can be force enabled by setting this value to true.
	// This means the config "enable" value is ignored for this middleware.
	// Only use this if the middleware is absolutely necessary.
	ForceEnable bool

	// Handler is the http.Handler of the middleware.
	Handler func(http.Handler) (handler http.Handler)
}

// middlewareRegistry stores all registered middleware, in order of registration.
// Middlewares are applied in the same order they were added.
var middlewareRegistry []Middleware

// RegisterMiddleware registers a new global middleware.
// Usually called from init() inside a middleware package.
func RegisterMiddleware(middleware Middleware) {
	middlewareRegistry = append(middlewareRegistry, middleware)
}

// ApplyRegisteredMiddleware wraps the given handler with all registered
// middleware functions, ordered by priority.
func applyRegisteredMiddleware(h http.Handler) http.Handler {
	// Sort the registered middlewares by priority
	sort.Slice(middlewareRegistry, func(i, j int) bool {
		return middlewareRegistry[i].Priority < middlewareRegistry[j].Priority
	})

	for _, middleware := range middlewareRegistry {
		var enabled bool = false

		// Determine if a middleware should be enabled.
		// ForceEnable skips the config.
		if middleware.ForceEnable {
			enabled = true
		} else {
			// Get enable value from config.
			config, okConfig := GetMiddlewareConfig(middleware.Name)
			if okConfig {
				enable_conf, okEnable := cfg.GetNestedValue[bool](config, "enable")
				if okEnable {
					enabled = enable_conf
				} else {
					// The config has no enable value.
					log.Printf("[MIDDLEWARE] %s has no enable value in config!", middleware.Name)
				}
			} else {
				// The module has no config entry.
				log.Printf("[MIDDLEWARE] %s has no config!", middleware.Name)
			}
		}

		if enabled {
			// Enable the middleware.
			h = middleware.Handler(h)
			log.Printf("[MIDDLEWARE] %s was applied!", middleware.Name)
		} else {
			// Middleware is disabled.
			log.Printf("[MIDDLEWARE] %s is disabled.", middleware.Name)
		}
	}
	return h
}
