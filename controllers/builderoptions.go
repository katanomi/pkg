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
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DefaultMaxConcurrentReconciles is the default number of max concurrent reconciles
const DefaultMaxConcurrentReconciles = 10

// BuilderOptionFunc is a function that can be used to configure the controller builder options
type BuilderOptionFunc func(options controller.Options) controller.Options

// BuilderOptions returns a functional set of options with conservative
// defaults.
func BuilderOptions(opts ...BuilderOptionFunc) controller.Options {
	options := DefaultOptions()
	for _, opt := range opts {
		options = opt(options)
	}
	return options
}

// DefaultOptions returns the default options for the controller
func DefaultOptions() controller.Options {
	return controller.Options{
		MaxConcurrentReconciles: DefaultMaxConcurrentReconciles,
		RateLimiter:             DefaultTypedRateLimiter[reconcile.Request](),
	}
}

// MaxConCurrentReconciles sets the max concurrent reconciles
func MaxConCurrentReconciles(num int) BuilderOptionFunc {
	return func(options controller.Options) controller.Options {
		options.MaxConcurrentReconciles = num
		return options
	}
}

// RateLimiter sets the rate limiter
func RateLimiter(rl workqueue.TypedRateLimiter[reconcile.Request]) BuilderOptionFunc {
	return func(options controller.Options) controller.Options {
		options.RateLimiter = rl
		return options
	}
}
