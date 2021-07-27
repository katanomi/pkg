/*
Copyright 2021 The Katanomi Authors.

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

type clientmanagerCtxKey struct{}

//WithManager sets a manager instance into a context
func WithManager(ctx context.Context, mgr *Manager) context.Context {
	return context.WithValue(ctx, clientmanagerCtxKey{}, mgr)
}

// ManagerCtx returns a *Manager in a given context. Returns nil if not found
func ManagerCtx(ctx context.Context) *Manager {
	val := ctx.Value(clientmanagerCtxKey{})
	if val == nil {
		return nil
	}
	return val.(*Manager)
}
