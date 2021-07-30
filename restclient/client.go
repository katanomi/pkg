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

package restclient

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type clientCtxKey struct{}

// WithRESTClient sets a client instance into a context
func WithRESTClient(ctx context.Context, clt *resty.Client) context.Context {
	return context.WithValue(ctx, clientCtxKey{}, clt)
}

// RESTClient returns a client.Client in a given context. Returns nil if not found
func RESTClient(ctx context.Context) *resty.Client {
	val := ctx.Value(clientCtxKey{})
	if val == nil {
		return nil
	}
	return val.(*resty.Client)
}
