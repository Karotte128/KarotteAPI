package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	cfg "github.com/karotte128/karottelib/config"
)

// InitAPI starts the HTTP server, loads all registered modules and middleware,
// and mounts each module under its prefix.
func RunAPI(ctx context.Context, config Config) error {
	// Load config
	loadConfig(config)

	// Get server config
	serverConfig, serverConfigOk := getServerConfig()
	if !serverConfigOk {
		return errors.New("[SERVER] No server config!")
	}

	// Get server address
	addr, addrOk := cfg.GetNestedValue[string](serverConfig, "address")
	if !addrOk {
		return errors.New("[SERVER] No server address config!")
	}

	if addr == "" {
		return errors.New("[SERVER] address is not configured!")
	}

	// Get ignoreStartupErrors
	ignoreStartupErrors, iseOk := cfg.GetNestedValue[bool](serverConfig, "ignoreStartupErrors")
	if !iseOk {
		return errors.New("[SERVER] No server ignoreStartupErrors config!")
	}

	// A multiplexer to route module-specific handlers.
	mux := http.NewServeMux()

	// Load all modules of the module registry.
	err := loadRegisteredModules(mux, ignoreStartupErrors)
	if err != nil {
		return err
	}

	// Apply global middlewares to the root mux.
	handler, err := applyRegisteredMiddlewares(mux, ignoreStartupErrors)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Channel for the HTTP server to report errors.
	serverErr := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("[SERVER] running on %s", addr)

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Wait for either:
	// 1. The server to fail
	// 2. The parent context to be cancelled
	select {
	case err := <-serverErr:
		return fmt.Errorf("[SERVER] server error: %w", err)

	case <-ctx.Done():
		log.Println("[SERVER] shutting down...")

		if err := server.Shutdown(context.Background()); err != nil {
			return err
		}

		shutdownRegisteredModules()

		log.Println("[SERVER] Done shutting down!")

		return nil
	}
}
