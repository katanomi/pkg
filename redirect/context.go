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

package redirect

import (
	"context"
)

type redirectKey struct{}

// WithRedirecter adds a redirect client to the context
func WithRedirecter(ctx context.Context, redi Interface) context.Context {
	return context.WithValue(ctx, redirectKey{}, redi)
}

// Redirecter returns a Redirecter in context
func Redirecter(ctx context.Context) Interface {
	val := ctx.Value(redirectKey{})
	if val == nil {
		return nil
	}
	clt, _ := val.(Interface)
	return clt
}