# KarotteAPI

KarotteAPI is a modular, extensible Go API framework designed to simplify building REST-style services. It provides shared types, middleware and module registration, and configuration utilities to help developers compose APIs in a structured and consistent way.

This README is intended for **developers who want to use KarotteAPI as a dependency** in their own projects.

---

## Table of Contents

- Overview
- Features
- Requirements
- Installation
- Core Packages
  - api
  - builtins
- Basic Usage
  - Setting up the API server
  - Registering a Module
  - Registering a Middleware
- Authentication
- Configuration
- License

---

## Overview

KarotteAPI is a Go module that defines common API abstractions (modules, middleware, permissions, and authentication context) and exposes a small public surface area for registering and interacting with them.

It is designed to be embedded into a larger application rather than run as a standalone server.

More utilities for setting up and using the API can be found in the [Karotte128/APIUtils](https://github.com/karotte128/apiutils) repository.
This package contains:
- Database helpers
- Config loader
- Permission provider (using the database)

This package is not a dependency, but its usage is highly recommended.

---

## Features

- Shared API interfaces for modules and middleware
- Centralized registration of API modules
- Permission and authentication abstractions
- Clean separation between public API and internal implementation

---

## Requirements

- Go 1.25 or newer
- Go modules enabled

---

## Installation

Add KarotteAPI as a dependency in your project:

```bash
go get github.com/Karotte128/KarotteAPI
```

Then import the required packages in your code:

```go
import (
    _ "github.com/Karotte128/KarotteAPI/builtins" // Enable builtin middlewares and modules. Highly recommended!
    "github.com/Karotte128/KarotteAPI/api"
)
```

---

## Core Packages

### `api`

The `api` package contains the supported API surface for interacting with the framework.

Types:

- `Config`
  Contains the API config as `map[string]any`.

- `Module`
  Represents a logical API module that can register routes, handlers, or behavior.

- `Middleware`
  Interface for request/response middleware components.

- `ModuleStatus`
  ModuleStatus stores the current status of all registered modules.

Functions:

- `InitAPI(api.Config)`
  Function used to set up and start the API server.

- `RegisterModule(module)`  
  Registers an API module.

- `GetModuleStatus()`
  Reads the status of all registered modules.

- `RegisterMiddleware(middleware)`  
  Registers a middleware component.

- `GetModuleConfig(moduleName)`  
  Retrieves configuration scoped to a specific module.

- `GetMiddlewareConfig(middlewareName)`  
  Retrieves configuration scoped to a specific middleware.

- `GetRequestContext[T any](r *http.Request, key string) (value T, ok bool)`
  Retrieves additional data from the request context.

- `SetRequestContext(r *http.Request, key string, value any)`
  Sets additional data on the request context.

All application-level interaction with KarotteAPI should go through this package.

### `builtins`

The builtins package contains usefull modules and middlewares.
To use these, import the module:
```go
import (
    _ "github.com/Karotte128/KarotteAPI/builtins"
)
```

#### Modules

- `health`
  Adds the /health endpoint, showing module status.

####

- `contenttype`
  Automatically sets the `Content-Type` header to `application/json` if it is not set.

- `logging`
  Adds basic request logging.

---

## Basic Usage

### Setting up the API server

To set up the API server, the `api.InitApi(api.Config)` function needs to be called with the `Config` as argument.

The following example shows a simple setup using [Karotte128/APIUtils](https://github.com/karotte128/apiutils).

```go
package main

import (
	"log"
    "context"

	"github.com/karotte128/apiutils/config"
	"github.com/karotte128/karotteapi/api"
	_ "github.com/karotte128/karotteapi/builtins"
)

func main() {
	err, rawConf := config.ReadConfigFromFile("config.toml") // Load toml config using APIUtils.
	if err != nil {
		log.Fatal("failed loading config: " + err.Error())
	}

	conf := config.ExpandEnvConfig(rawConf) // Replace ENV vars in the config (APIUtils)

	api.InitAPI(conf) // Start the API server.
}
```

The config.toml for this example contains:

```toml
[server]
address = "${ADDR:-:8080}"

[modules.health]
enable = true
```

---

### Registering a Module

Modules are typically registered during initialization:

```go
var exampleModule = api.Module{ // Create the module info.
	Name:     "example", // Name of the module, used for logging.
	Routes:   routes, // Function that provides the API routes of the module.
	Startup:  startup, // Function that is executed after the module has been registered. Can be nil if not needed.
	Shutdown: shutdown, // Function that is executed before the server shuts down. Can be nil if not needed.
}

func routes() (string, http.Handler) {
	mux := http.NewServeMux() // Create the new http mux.
	mux.HandleFunc("/example", example) // Add handler func.
	return "/example", mux // Return prefix and mux.
}

func startup() error {
	log.Println("[MODULE] Starting the example module!")
	return nil // Return nil (no error).
}

func shutdown() error {
	log.Println("[MODULE] Shutting down the example module!")
	return nil // Return nil (no error).
}

func init() { // init() is used to register the module before the server starts.
	api.RegisterModule(statusModule) // Add the module to the registry.
}

func example(w http.ResponseWriter, r *http.Request) { // http handler that handles the request.
    fmt.Fprint(w, "Hello World!") // Write the http response.
}
```

Your module must implement the appropriate `api.Module` interface.

---

### Registering a Middleware

```go
var exampleMiddleware = api.Middleware{ // Create the example middleware.
	Name:     "example", // Name of the middleware, used for logging.
	Handler:  exampleHandler, // Handler function to modify the request.
	Priority: 10, // Higher number -> gets applied later; lower number -> gets applied earlier.
	ForceEnable: false, // If enabled, the enable value in the middleware config is ignored. Only use on necessary middlewares.
}

func exampleHandler(next http.Handler) http.Handler { // Handler function, returns the new (modified) handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Return the new handler

        w.Header().Set("Example-Header", "Hello World!") // Set a header (example).

		next.ServeHTTP(w, r) // Serve the next handler in the chain.
	})
}

func init() { // init() is used to register the middleware before the server starts.
	api.RegisterMiddleware(exampleMiddleware) // Add the middleware to the registry
}
```

Middleware can inspect or modify requests using the provided context.

---

## Authentication

KarotteAPI does not have authentication built in.
For an easy auth system it is recommended to use [Karotte128/APIUtils simpleauth](https://github.com/karotte128/apiutils/tree/main/simpleauth)

## Configuration

KarotteAPI expects configuration to be supplied by the host application.  
The mechanism (environment variables, files, flags, etc.) is left to the integrator.

[Karotte128/APIUtils](https://github.com/karotte128/apiutils) contains a simple to use configuration loader system that is compatible with this API.
It reads the config from a `.toml` file and replaces `${ENV}` variables dynamically.

Use `api.GetModuleConfig(name)` inside a module to get module-specific configuration.

---

## License

This project is licensed under the MIT License.
