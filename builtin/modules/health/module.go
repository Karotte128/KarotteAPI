package health

import (
	"log"
	"net/http"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/core"
)

var healthModule = karotteapi.Module{
	Name:     "health",
	Routes:   routes,
	Startup:  startup,
	Shutdown: shutdown,
}

func routes() (string, http.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	return "/health", mux
}

func startup() error {
	log.Println("[MODULE] Starting the health module!")
	return nil
}

func shutdown() error {
	log.Println("[MODULE] Shutting down the health module!")
	return nil
}

func init() {
	core.RegisterModule(healthModule)
}
