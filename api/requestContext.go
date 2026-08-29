package api

import (
	"context"
	"net/http"
)

// This adds the info to the request data.
// It is usually used by a middleware.
func SetRequestContext(r *http.Request, key string, value any) {
	newCtx := context.WithValue(r.Context(), key, value)
	*r = *r.WithContext(newCtx)
}

// This retrieves the info from the request context.
// It is usually used in a module.
func GetRequestContext[T any](r *http.Request, key string) (value T, ok bool) {
	v := r.Context().Value(key)
	if v == nil {
		// Context key not found
		var zero T
		return zero, false
	}

	// Attempt a type‑assert to the caller‑requested generic type.
	typed, castOk := v.(T)
	return typed, castOk
}
