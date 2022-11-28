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

// Package restclient contains functions to add and retrieve rest client from context
package restclient

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type clientCtxKey struct{}

// WithRESTClient sets a https client instance into a context
func WithRESTClient(ctx context.Context, clt *resty.Client) context.Context {
	return context.WithValue(ctx, clientCtxKey{}, clt)
}

// RESTClient returns a https client.Client in a given context. Returns nil if not found
func RESTClient(ctx context.Context) *resty.Client {
	val := ctx.Value(clientCtxKey{})
	if val == nil {
		return nil
	}
	return val.(*resty.Client)
}

// type httpClientCtxKey struct{}

// // WithHttpRESTClient sets a http client instance into a context
// func WithHttpRESTClient(ctx context.Context, clt *resty.Client) context.Context {
// 	return context.WithValue(ctx, httpClientCtxKey{}, clt)
// }

// // HttpRESTClient returns a http client.Client in a given context. Returns nil if not found
// func HttpRESTClient(ctx context.Context) *resty.Client {
// 	val := ctx.Value(httpClientCtxKey{})
// 	if val == nil {
// 		return nil
// 	}
// 	return val.(*resty.Client)
// }
