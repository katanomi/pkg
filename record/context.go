/*
Copyright 2024 The AlaudaDevops Authors.

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

// Package record has all client logic for recording and reporting
package record

import (
	"context"

	"k8s.io/client-go/tools/record"
)

type recordCtxKey struct{}

// WithRecorder adds a recorder to the context
func WithRecorder(ctx context.Context, r record.EventRecorder) context.Context {
	return context.WithValue(ctx, recordCtxKey{}, r)
}

// FromContext gets a record from the context
func FromContext(ctx context.Context) record.EventRecorder {
	r, _ := ctx.Value(recordCtxKey{}).(record.EventRecorder)
	return r
}
