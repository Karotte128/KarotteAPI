package internal

import (
	"context"

	"github.com/karotte128/karotteapi"
)

// This adds the info to the request data.
// It is usually used by a middleware.
func SetRequestContext(ctx context.Context, info *karotteapi.RequestContext) context.Context {
	return context.WithValue(ctx, info.ContextKey, info.Info)
}

// This retrieves the info from the request context.
// It is usually used in a module.
func GetRequestContext(ctx context.Context, contextKey string) karotteapi.RequestContext {
	requestContext := karotteapi.RequestContext{
		Info:       ctx.Value(contextKey),
		ContextKey: contextKey,
	}
	return requestContext
}
