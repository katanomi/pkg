package scheme

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type schemeCtxKey struct{}

func WithScheme(ctx context.Context, scheme *runtime.Scheme) context.Context {
	return context.WithValue(ctx, schemeCtxKey{}, scheme)
}

func Scheme(ctx context.Context) *runtime.Scheme {
	val := ctx.Value(schemeCtxKey{})
	if val == nil {
		return nil
	}
	return val.(*runtime.Scheme)
}
