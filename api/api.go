package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/karotte128/karotteapi"
	_ "github.com/karotte128/karotteapi/builtin/middleware" // automatically loads all middleware via init()
	_ "github.com/karotte128/karotteapi/builtin/modules"    // automatically loads all modules via init()
	"github.com/karotte128/karotteapi/internal"
)

// InitAPI starts the HTTP server, loads all registered modules and middleware,
// and mounts each module under its prefix.
func InitAPI(details karotteapi.ApiDetails) {
	// Load config
	internal.LoadConfig(details.Config)

	// Get server config
	serverConfig, serverConfigOk := internal.GetServerConfig()
	if !serverConfigOk {
		log.Fatal("[SERVER] No server config!")
	}

	// Get server address
	addr, addrOk := internal.GetNestedValue[string](serverConfig, "address")
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
	internal.LoadRegisteredModules(mux)

	// Apply global middleware to the root mux.
	handler := internal.ApplyRegisteredMiddleware(mux)

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
	internal.ShutdownRegisteredModules()
}
