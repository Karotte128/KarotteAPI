package api

import (
	"errors"
	"fmt"
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

	// Startup is a function that is run on startup.
	// This can be used to initialize a connection to external services like databases.
	Startup func() error

	// Shutdown is a function that is run on shutdown.
	// This can be used to cleanly disconnect from services connected during startup.
	Shutdown func() error
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
func applyRegisteredMiddlewares(h http.Handler, ignoreErrors bool) (http.Handler, error) {
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
			// Attempt to load middleware config.
			config, okConfig := GetMiddlewareConfig(middleware.Name)
			if !okConfig {
				fErr := fmt.Sprintf("[MIDDLEWARE] %s has no config!", middleware.Name)
				log.Println(fErr)

				if !ignoreErrors {
					return nil, errors.New(fErr)
				}

				break
			}

			// Get enable value from config.
			enable_conf, okEnable := cfg.GetNestedValue[bool](config, "enable")
			if !okEnable {
				fErr := fmt.Sprintf("[MIDDLEWARE] %s has no enable value in config!", middleware.Name)
				log.Println(fErr)

				if !ignoreErrors {
					return nil, errors.New(fErr)
				}

				break
			}

			enabled = enable_conf
		}

		// Skip further processing if middleware is disabled.
		if !enabled {
			log.Printf("[MIDDLEWARE] %s is disabled.", middleware.Name)
			break
		}

		err := safeStartmiddleware(middleware)
		if err != nil {
			log.Println(err)

			if !ignoreErrors || middleware.ForceEnable {
				return nil, err
			}

			break
		}

		// Enable the middleware.
		h = middleware.Handler(h)
		log.Printf("[MIDDLEWARE] %s was applied!", middleware.Name)
	}

	return h, nil
}

// ShutdownRegisteredMiddlewares shuts down all middlewares.
func shutdownRegisteredMiddlewares() {
	for _, middleware := range middlewareRegistry {
		err := safeShutdownMiddleware(middleware)
		if err != nil {
			// Error or panic occured while trying to shutdown middleware
			log.Println(err)
			break
		}
	}
}

// safeShutdownMiddleware is a function that attempts to execute the shutdown function of a middleware.
// It returns nil if the shutdown is successfull or the middleware does not provide a shutdown function.
// It makes sure that a panic in the shutdown function does not crash the server.
func safeShutdownMiddleware(middleware Middleware) (err error) {
	// return immediately if shutdown is not needed.
	if middleware.Shutdown == nil {
		return nil
	}

	// recover panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[MIDDLEWARE] %s panicked during shutdown: %v", middleware.Name, r)
		}
	}()

	if startupErr := middleware.Shutdown(); startupErr != nil {
		return fmt.Errorf("[MIDDLEWARE] %s failed shutdown: %w", middleware.Name, startupErr)
	}

	return nil
}

// safeStartMiddleware is a function that attempts to execute the startup function of a middleware.
// It returns nil if the startup is successfull or the middleware does not provide a startup function.
// It makes sure that a panic in the startup function does not crash the server.
func safeStartmiddleware(middleware Middleware) (err error) {
	// return immediately if startup is not needed.
	if middleware.Startup == nil {
		return nil
	}

	// recover panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[MIDDLEWARE] %s panicked during startup: %v", middleware.Name, r)
		}
	}()

	if startupErr := middleware.Startup(); startupErr != nil {
		return fmt.Errorf("[MIDDLEWARE] %s failed startup: %w", middleware.Name, startupErr)
	}

	return nil
}
