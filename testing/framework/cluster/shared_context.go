/*
Copyright 2023 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"context"

	kclient "github.com/AlaudaDevops/pkg/client"
	"github.com/AlaudaDevops/pkg/testing/framework/base"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/injection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type sharedSchemeCtxKey struct{}

// WithSharedScheme wrap scheme into context
func WithSharedScheme(ctx context.Context, scheme *runtime.Scheme) context.Context {
	return context.WithValue(ctx, sharedSchemeCtxKey{}, scheme)
}

// FromSharedScheme get scheme from context
func FromSharedScheme(ctx context.Context) *runtime.Scheme {
	val := ctx.Value(sharedSchemeCtxKey{})
	if val == nil {
		return nil
	}
	return val.(*runtime.Scheme)
}

// ShareScheme to construct extension to share scheme
func ShareScheme(scheme *runtime.Scheme) base.SharedExtension {
	return base.SharedExtensionFunc(func(ctx context.Context) context.Context {
		return WithSharedScheme(ctx, scheme)
	})
}

// SharedClient to construct extension to share client
func SharedClient() base.SharedExtension {
	return base.SharedExtensionFunc(func(ctx context.Context) context.Context {
		cfg := injection.GetConfig(ctx)
		if cfg == nil {
			cfg = ctrl.GetConfigOrDie()
			ctx = injection.WithConfig(ctx, cfg)
		}

		scheme := FromSharedScheme(ctx)
		client, err := client.New(cfg, client.Options{Scheme: scheme})
		if err != nil {
			panic(err)
		}
		return kclient.WithDirectClient(ctx, client)
	})
}
