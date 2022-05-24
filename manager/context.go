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

// Package manager contains functions to add and retrieve manager from context
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
