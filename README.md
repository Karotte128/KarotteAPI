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
  - karotteapi
  - api
  - core
- Basic Usage
  - Setting up the API server
  - Registering a Module
  - Registering a Middleware
  - Permission Checks
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
    "github.com/Karotte128/KarotteAPI"
    "github.com/Karotte128/KarotteAPI/api"
    "github.com/Karotte128/KarotteAPI/core"
)
```

---

## Core Packages

### karotteapi

The `karotteapi` package defines the core interfaces and data structures used by the framework.

Key concepts include:

- **ApiDetails**
  Contains the API config as `map[string]any` and the `PermissionProvider`.

- **Module**  
  Represents a logical API module that can register routes, handlers, or behavior.

- **Middleware**  
  Interface for request/response middleware components.

- **PermissionProvider**  
  Abstraction for permission resolution and checks.

These types are intended to be implemented or consumed by your application code.

---

### api

The `api` package contains the `InitAPI` function used to set up and start the API server. 

---

### core

The `core` package exposes the supported API surface for interacting with the framework.

Common functions include:

- `RegisterModule(module)`  
  Registers an API module.

- `RegisterMiddleware(middleware)`  
  Registers a middleware component.

- `HasPermission(ctx, permission)`
  Checks whether a request has a specific permission from a request context.

- `GetModuleConfig(moduleName)`  
  Retrieves configuration scoped to a specific module.

All application-level interaction with KarotteAPI should go through this package.

---

## Basic Usage

### Setting up the API server

To set up the API server, the `api.InitApi()` function needs to be called with `ApiDetails` as argument.

The following example shows a simple setup using [Karotte128/APIUtils](https://github.com/karotte128/apiutils).

```go
package main

import (
	"log"
    "context"

	"github.com/karotte128/apiutils/config"
    "github.com/karotte128/apiutils/permissions"
	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/api"
	"github.com/karotte128/karotteapi/core"

    "github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool // pgx connection pool

func main() {
	var details karotteapi.ApiDetails // Create empty ApiDetails.

	err, rawConf := config.ReadConfigFromFile("config.toml") // Load toml config using APIUtils.
	if err != nil {
		log.Fatal("failed loading config: " + err.Error())
	}

	conf := config.ExpandEnvConfig(rawConf) // Replace ENV vars in the config (APIUtils)

	dbconn, ok := core.GetNestedValue[string](conf, "database", "connection") // Get the pgx connection string from the config.
	if !ok {
		log.Fatal("no database config!")
	}

	createConnection(dbconn) // Create database connection.

	details.Config = conf // Set the config.
	details.PermProvider = getPermissionWrapper // Set the permission provider.
	api.InitAPI(details) // Start the API server.
}

func createConnection(pgxconf string) { // Create database connection using the config value.
	poolConfig, err := pgxpool.ParseConfig(pgxconf)
	ConnPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = ConnPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping to database failed: %v\n", err)
	}
}

func getPermissionWrapper(key string) []string { // Wrapper around APIUtils GetPermission
	info, err := permissions.GetPermission(ConnPool, "authentication", key)
	if err != nil {
		log.Println(err.Error())
	}

	return info.Permissions
}
```

The config.toml for this example contains:

```toml
[server]
address = "${ADDR:-:8080}"

[modules.status]
enable = true

[database]
connection = "${DBCONN}"
```

---

### Registering a Module

Modules are typically registered during initialization:

```go
var exampleModule = karotteapi.Module{ // Create the module info.
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
	core.RegisterModule(statusModule) // Add the module to the registry.
}

func example(w http.ResponseWriter, r *http.Request) { // http handler that handles the request.
    fmt.Fprint(w, "Hello World!") // Write the http response.
}
```

Your module must implement the appropriate `karotteapi.Module` interface.

---

### Registering a Middleware

```go
var exampleMiddleware = karotteapi.Middleware{ // Create the example middleware.
	Name:     "example", // Name of the middleware, used for logging.
	Handler:  exampleHandler, // Handler function to modify the request.
	Priority: 10, // Higher number -> gets applied later; lower number -> gets applied earlier.
}

func exampleHandler(next http.Handler) http.Handler { // Handler function, returns the new (modified) handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Return the new handler

        w.Header().Set("Example-Header", "Hello World!") // Set a header (example).

		next.ServeHTTP(w, r) // Serve the next handler in the chain.
	})
}

func init() { // init() is used to register the middleware before the server starts.
	core.RegisterMiddleware(exampleMiddleware) // Add the middleware to the registry
}
```

Middleware can inspect or modify requests using the provided context.

---

### Permission Checks

Inside a module handler:

```go
func handle(w http.ResponseWriter, r *http.Request) {
    if !core.HasPermission(r.Context(), "admin:*") { // Get AuthInfo from the request context and check if user has permission 
        fmt.Fprint(w, "Access denied!") // User does not have permission
    }
    
    fmt.Fprint(w, "Hello World!") // User does have permission   
}
```

---

## Configuration

KarotteAPI expects configuration to be supplied by the host application.  
The mechanism (environment variables, files, flags, etc.) is left to the integrator.

[Karotte128/APIUtils](https://github.com/karotte128/apiutils) contains a simple to use configuration loader system that is compatible with this API.
It reads the config from a `.toml` file and replaces `${ENV}` variables dynamically.

Use `core.GetModuleConfig(name)` inside a module to get module-specific configuration.

---

## License

This project is licensed under the GNU General Public License v3.0.
