package karotteapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/karotte128/karotteapi/internal/core"
	_ "github.com/karotte128/karotteapi/internal/middleware" // automatically loads all middleware via init()
	_ "github.com/karotte128/karotteapi/internal/modules"    // automatically loads all modules via init()
)

// ApiDetails contains all details to create a new api.
type ApiDetails struct {
	ConfigPath   string
	PermProvider core.PermissionProvider
}

// InitAPI starts the HTTP server, loads all registered modules and middleware,
// and mounts each module under its prefix.
func InitAPI(details ApiDetails) {
	// Load config
	err := core.LoadConfig(details.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	// Get server config
	serverConfig, serverConfigOk := core.GetServerConfig()
	if !serverConfigOk {
		log.Fatal("[SERVER] No server config!")
	}

	// Get server address
	addr, addrOk := core.GetNestedValue[string](serverConfig, "address")
	if !addrOk {
		log.Fatal("[SERVER] No server address config!")
	}

	// check if address has no value
	if addr == "" {
		log.Fatal("[SERVER] address is not configured!")
	}

	// A multiplexer to route module-specific handlers.
	mux := http.NewServeMux()

	core.SetPermissionProvider(details.PermProvider)

	// Load all modules of the module registry.
	core.LoadRegisteredModules(mux)

	// Apply global middleware to the root mux.
	handler := core.ApplyRegisteredMiddleware(mux)

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
	core.ShutdownRegisteredModules()
}
