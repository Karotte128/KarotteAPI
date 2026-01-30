package middleware

import (
	"log"
	"net/http"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/core"
)

// RecoveryMiddleware wraps every request handler in a recover() block.
// If a module panics, the panic is caught here so the server remains alive.
// The user receives a 500 error and the panic is logged.

var recoveryMiddleware = karotteapi.Middleware{
	Name:     "recovery",
	Handler:  recoveryHandler,
	Priority: 0,
}

func recoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[RECOVERY] Panic caught: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Automatically register this middleware at startup.
// It becomes part of the global middleware registry.
func init() {
	core.RegisterMiddleware(recoveryMiddleware)
}
