package status

import (
	"log"
	"net/http"

	"github.com/karotte128/karotteapi/apitypes"
	"github.com/karotte128/karotteapi/internal/core"
)

var statusModule = apitypes.Module{
	Name:     "status",
	Routes:   routes,
	Startup:  startup,
	Shutdown: shutdown,
}

func routes() (string, http.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", status)
	return "/status", mux
}

func startup() error {
	log.Println("[MODULE] Starting the status module!")
	return nil
}

func shutdown() error {
	log.Println("[MODULE] Shutting down the status module!")
	return nil
}

func init() {
	core.RegisterModule(statusModule)
}
