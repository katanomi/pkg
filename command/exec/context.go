/*
Copyright 2022 The AlaudaDevops Authors.

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

package exec

import "context"

// cmderKey key for Cmder in Context
type cmderKey struct{}

// FromContextCmder returns a Cmder from the context
func FromContextCmder(ctx context.Context) (cmdr Cmder) {
	value := ctx.Value(cmderKey{})
	if value != nil {
		if cmdr, ok := value.(Cmder); ok {
			return cmdr
		}
	}
	return DefaultCmder
}

// WithCmder adds a Cmder to a context
func WithCmder(ctx context.Context, cmdr Cmder) context.Context {
	return context.WithValue(ctx, cmderKey{}, cmdr)
}

// cmdKey key for Cmd in Context
type cmdKey struct{}

// FromContextCmd returns a Cmd from the context
func FromContextCmd(ctx context.Context) (cmdr Cmd) {
	if value := ctx.Value(cmdKey{}); value != nil {
		if cmdr, ok := value.(Cmd); ok {
			return cmdr
		}
	}
	return nil
}

// WithCmd adds a Cmd to a context
func WithCmd(ctx context.Context, cmdr Cmd) context.Context {
	return context.WithValue(ctx, cmdKey{}, cmdr)
}
