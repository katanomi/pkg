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

package controllers

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
)

type rateLimiterKey struct{}

// WithRateLimiter stores a RateLimiter into context
func WithRateLimiter(ctx context.Context, rl workqueue.RateLimiter) context.Context {
	return context.WithValue(ctx, rateLimiterKey{}, rl)
}

// RateLimiterCtx retrieves a RateLimiter from context. Returns nil if none
func RateLimiterCtx(ctx context.Context) workqueue.RateLimiter {
	val := ctx.Value(rateLimiterKey{})
	if val == nil {
		return nil
	}
	return val.(workqueue.RateLimiter)
}

type requestKey struct{}

// WithReconcileRequest stores a Request key into context
func WithReconcileRequest(ctx context.Context, key ctrl.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, key)
}

// ReconcileRequestCtx retrieves a Request key from context
// returns an empty object if none
func ReconcileRequestCtx(ctx context.Context) ctrl.Request {
	val := ctx.Value(requestKey{})
	if val == nil {
		return ctrl.Request{}
	}
	return val.(ctrl.Request)
}

// numRequeuesCtxKey returns back how many failures the item has had
type numRequeuesCtxKey struct{}

// WithNumRequeues adds the numRequeues to the context
func WithNumRequeues(ctx context.Context, num int) context.Context {
	return context.WithValue(ctx, numRequeuesCtxKey{}, num)
}

// GetNumRequeues gets numRequeue from the context
func GetNumRequeues(ctx context.Context) int {
	num, _ := ctx.Value(numRequeuesCtxKey{}).(int)
	return num
}
