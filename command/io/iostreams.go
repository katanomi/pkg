/*
Copyright 2022 The Katanomi Authors.

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

package io

import (
	"context"
	"os"

	clioptions "k8s.io/cli-runtime/pkg/genericclioptions"
)

// key for reading/writing into context
type ioStreamsKey struct{}

// WithIOStreams adds IOStreams into the context
func WithIOStreams(ctx context.Context, ioStreams *clioptions.IOStreams) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ioStreamsKey{}, ioStreams)
}

// GetIOStreams returns IOStreams stored in the context if any
// if not found will return nil *IOStreams
func GetIOStreams(ctx context.Context) (ioStreams *clioptions.IOStreams) {
	if ctx == nil {
		ctx = context.Background()
	}
	if val := ctx.Value(ioStreamsKey{}); val != nil {
		ioStreams = val.(*clioptions.IOStreams)
	}
	return
}

// MustGetIOStreams gets the IOStream from context or initiates a default
// using os.Stdin, os.Stout, os.Sterr
func MustGetIOStreams(ctx context.Context) (ioStreams *clioptions.IOStreams) {
	if ioStreams = GetIOStreams(ctx); ioStreams == nil {
		ioStreams = &clioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		}
	}
	return
}
