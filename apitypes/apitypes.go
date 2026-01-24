package apitypes

import "net/http"

// ApiDetails contains all details to create a new api.
type ApiDetails struct {
	Config       map[string]any
	PermProvider PermissionProvider
}

// AuthInfo is created by the auth middleware.
// It contains the authentication status and permissions of the request.
type AuthInfo struct {
	// ApiKey is the raw key sent by the user. Do not use this.
	ApiKey string

	// Permissions is the list of permissions the user has.
	Permissions []string
}

// PermissionProvider is the function used for getting the users permissions from the API key.
type PermissionProvider func(key string) []string

// Middleware is the struct the middleware needs to provide to the middleware registry to register itself.
type Middleware struct {
	// Name is the name of the middleware. It is used for logging.
	Name string

	// Priority is the order in which middlewares should be registered.
	// Lower number means the middleware gets registered earlier (higher priority).
	Priority uint

	// Handler is the http.Handler of the middleware.
	Handler func(http.Handler) (handler http.Handler)
}

// Module is the struct the module needs to provide to the module registry to register itself.
type Module struct {
	// Name is the name of the module. It is used for logging.
	Name string

	// Routes returns a URL prefix and an http.Handler that serves all routes
	// for this module.
	//
	// Example:
	//   prefix = "/example/"
	//   handler = http.HandlerFunc()
	Routes func() (prefix string, handler http.Handler)

	// Startup is a function that is run on startup.
	// This can be used to initialize a connection to external services like databases.
	Startup func() error

	// Shutdown is a function that is run on shutdown.
	// This can be used to cleanly disconnect from services connected during Startup().
	Shutdown func() error
}
