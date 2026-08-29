package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/karotte128/karottelib/config"
)

// InitAPI starts the HTTP server, loads all registered modules and middleware,
// and mounts each module under its prefix.
func InitAPI(config Config) {
	// Load config
	loadConfig(config)

	// Get server config
	serverConfig, serverConfigOk := getServerConfig()
	if !serverConfigOk {
		log.Fatal("[SERVER] No server config!")
	}

	// Get server address
	addr, addrOk := cfg.GetNestedValue[string](serverConfig, "address")
	if !addrOk {
		log.Fatal("[SERVER] No server address config!")
	}

	// check if address has no value
	if addr == "" {
		log.Fatal("[SERVER] address is not configured!")
	}

	// A multiplexer to route module-specific handlers.
	mux := http.NewServeMux()

	// Load all modules of the module registry.
	loadRegisteredModules(mux)

	// Apply global middleware to the root mux.
	handler := applyRegisteredMiddleware(mux)

	// listen for shutdown notification
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// start http server
	go func() {
		log.Printf("[SERVER] running on %s", addr)
		err := http.ListenAndServe(addr, handler)

		if err != nil {
			log.Fatalf("[SERVER] error: %v", err)
		}
	}()

	// shutdown triggered
	<-ctx.Done()

	log.Println("[SERVER] shutting down...")
	// shutting down registered modules
	shutdownRegisteredModules()
}
