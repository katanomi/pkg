package manager

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type managerCtxKey struct{}

// WithManager sets a manager instance into a context
func WithManager(ctx context.Context, mgr manager.Manager) context.Context {
	return context.WithValue(ctx, managerCtxKey{}, mgr)
}

// Manager returns a manager in a given context. Returns nil if not found
func Manager(ctx context.Context) manager.Manager {
	val := ctx.Value(managerCtxKey{})
	if val == nil {
		return nil
	}
	return val.(manager.Manager)
}
