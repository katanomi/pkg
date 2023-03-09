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

// Package client contains functions to add and retrieve client from context
package client

import (
	"context"
	"fmt"

	"k8s.io/apiserver/pkg/endpoints/request"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
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

type directClientCtxKey struct{}

// WithDirectClient sets a client instance into a context
func WithDirectClient(ctx context.Context, clt client.Client) context.Context {
	return context.WithValue(ctx, directClientCtxKey{}, clt)
}

// DirectClient returns a client.Client in a given context. Returns nil if not found
func DirectClient(ctx context.Context) client.Client {
	val := ctx.Value(directClientCtxKey{})
	if val == nil {
		return nil
	}
	return val.(client.Client)
}

type clientmanagerCtxKey struct{}

// WithManager sets a manager instance into a context
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

type dynamicClientCtxKey struct{}

// WithDynamicClient sets a dynamic.Interface client instance into a context
func WithDynamicClient(ctx context.Context, client dynamic.Interface) context.Context {
	return context.WithValue(ctx, dynamicClientCtxKey{}, client)
}

// DynamicClient returns a dynamic client.Client, returns nil if not found
func DynamicClient(ctx context.Context) (dynamic.Interface, error) {
	val := ctx.Value(dynamicClientCtxKey{})
	if val == nil {
		return nil, fmt.Errorf("not found")
	}

	return val.(dynamic.Interface), nil
}

type clusterCtxKey struct{}

// WithCluster sets a cluster.Cluster instance into a context
func WithCluster(ctx context.Context, client cluster.Cluster) context.Context {
	return context.WithValue(ctx, clusterCtxKey{}, client)
}

// Cluster returns a cluster.Cluster, returns nil if not found
func Cluster(ctx context.Context) cluster.Cluster {
	val := ctx.Value(clusterCtxKey{})
	if val == nil {
		return nil
	}
	return val.(cluster.Cluster)
}

// User returns a user.Info, returns nil if not found
func User(ctx context.Context) user.Info {
	u, _ := request.UserFrom(ctx)
	return u
}

// cfgKeyOfApp is the key that the config make is associated with.
type cfgKeyOfApp struct{}

// WithAppConfig associates a given config with the app context.
func WithAppConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, cfgKeyOfApp{}, cfg)
}

// GetAppConfig gets the current config of app (pod) from the context.
func GetAppConfig(ctx context.Context) *rest.Config {
	value := ctx.Value(cfgKeyOfApp{})
	if value == nil {
		return nil
	}
	return value.(*rest.Config)
}
