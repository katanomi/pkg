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
)

type pluginClientKey struct{}

// WithPluginClient returns a copy of parent in which the pluginClient value is set
func WithPluginClient(parent context.Context, pluginClient *PluginClient) context.Context {
	return context.WithValue(parent, pluginClientKey{}, pluginClient)
}

// PluginClientFrom returns the value of the pluginClient key on the ctx
func PluginClientFrom(ctx context.Context) (*PluginClient, bool) {
	pluginClient, ok := ctx.Value(pluginClientKey{}).(*PluginClient)
	return pluginClient, ok
}

// PluginClientValue returns the value of the pluginClient key on the ctx, or the nil if none
func PluginClientValue(ctx context.Context) *PluginClient {
	pluginClient, _ := PluginClientFrom(ctx)
	return pluginClient
}
