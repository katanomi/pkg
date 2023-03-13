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

package storage

import "context"

type storagePluginNameKey struct{}

type clientSetKey struct{}

// CtxWithPluginName return context.Context with plugin name
func CtxWithPluginName(ctx context.Context, pluginName string) context.Context {
	return context.WithValue(ctx, storagePluginNameKey{}, pluginName)
}

// PluginNameFromCtx returns plugin name from context
func PluginNameFromCtx(ctx context.Context) string {
	val := ctx.Value(storagePluginNameKey{})
	if val == nil {
		return ""
	}
	return val.(string)
}

// WithClientSet return context with client set
func WithClientSet(ctx context.Context, clientSet Interface) context.Context {
	return context.WithValue(ctx, clientSetKey{}, clientSet)
}

// ClientSetFromCtx return clientset from context
func ClientSetFromCtx(ctx context.Context) Interface {
	val := ctx.Value(clientSetKey{})
	if val == nil {
		return nil
	}
	return val.(Interface)
}
