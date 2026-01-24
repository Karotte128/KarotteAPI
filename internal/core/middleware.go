package core

import (
	"log"
	"net/http"
	"sort"

	"github.com/karotte128/karotteapi/apitypes"
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

// registry stores all registered middleware, in order of registration.
// Middlewares are applied in the same order they were added.
var middleware_registry []apitypes.Middleware

// RegisterMiddleware registers a new global middleware.
// Usually called from init() inside a middleware package.
func RegisterMiddleware(middleware apitypes.Middleware) {
	middleware_registry = append(middleware_registry, middleware)
}

// Middlewares returns all registered middleware.
func GetMiddlewares() []apitypes.Middleware {
	return middleware_registry
}

// ApplyRegisteredMiddleware wraps the given handler with all registered
// middleware functions in registration order.
func ApplyRegisteredMiddleware(h http.Handler) http.Handler {
	sort.Slice(middleware_registry, func(i, j int) bool {
		return middleware_registry[i].Priority < middleware_registry[j].Priority
	})

	for _, middleware := range middleware_registry {
		h = middleware.Handler(h)
		log.Printf("[MIDDLEWARE] %s was applied!", middleware.Name)
	}
	return h
}
