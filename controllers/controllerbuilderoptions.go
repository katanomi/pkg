/*
Copyright 2023 The AlaudaDevops Authors.

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

package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

// ControllerBuilderOption is a function that can be used to configure the controller builder
type ControllerBuilderOption func(*builder.Builder) *builder.Builder

// ApplyControllerBuilderOptions applies the given options to the controller builder
func ApplyControllerBuilderOptions(builder *builder.Builder, funcs ...ControllerBuilderOption) *builder.Builder {
	for _, builderFunc := range funcs {
		builder = builderFunc(builder)
	}
	return builder
}

// WithBuilderOptions returns a ControllerBuilderOption that applies the given options to the controller builder
func WithBuilderOptions(opts controller.Options) ControllerBuilderOption {
	return func(b *builder.Builder) *builder.Builder {
		b = b.WithOptions(opts)
		return b
	}
}
