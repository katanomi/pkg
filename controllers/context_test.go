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
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestRateLimiterContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	limiter := RateLimiterCtx(ctx)
	g.Expect(limiter).To(BeNil())

	limiter = DefaultRateLimiter()
	ctx = WithRateLimiter(ctx, limiter)
	g.Expect(RateLimiterCtx(ctx)).To(Equal(limiter))
}

func TestManagerContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	req := ReconcileRequestCtx(ctx)
	g.Expect(req).To(Equal(ctrl.Request{}))

	req = ctrl.Request{NamespacedName: types.NamespacedName{Name: "abc", Namespace: "default"}}
	ctx = WithReconcileRequest(ctx, req)
	g.Expect(ReconcileRequestCtx(ctx)).To(Equal(req))
}
