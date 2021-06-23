package client

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type clientCtxKey struct{}

// WithClient sets a client instance into a context
func WithClient(ctx context.Context, clt client.Client) context.Context {
	return context.WithValue(ctx, clientCtxKey{}, clt)
}

// Client returns a client.Client in a given context. Returns nil if not found
func Client(ctx context.Context) client.Client {
	val := ctx.Value(clientCtxKey{})
	if val == nil {
		return nil
	}
	return val.(client.Client)
}
