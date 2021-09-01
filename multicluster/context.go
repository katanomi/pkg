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

package multicluster

import (
	"context"
)

type multiClusterKey struct{}

// WithMultiCluster adds a multi cluster client to the context
func WithMultiCluster(ctx context.Context, clt Interface) context.Context {
	return context.WithValue(ctx, multiClusterKey{}, clt)
}

// MultiCluster returns a multicluster client in context
func MultiCluster(ctx context.Context) Interface {
	val := ctx.Value(multiClusterKey{})
	if val == nil {
		return nil
	}
	clt, _ := val.(Interface)
	return clt
}
