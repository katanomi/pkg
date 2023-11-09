/*
Copyright 2023 The Katanomi Authors.

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
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestNewScheduleReconciler(t *testing.T) {
	r := NewScheduleReconciler(nil)
	sr, ok := r.(*scheduleReconciler)

	g := NewGomegaWithT(t)
	g.Expect(ok).To(BeTrue())
	g.Expect(sr.syncPeriod).To(Equal(minSyncPeriod))

	r = NewScheduleReconciler(reconcile.Func(emptyReconciler), func(options *ScheduleOptions) {
		options.SyncPeriod = 1 * time.Minute
	})
	sr, ok = r.(*scheduleReconciler)
	g.Expect(ok).To(BeTrue())
	g.Expect(sr.syncPeriod).To(Equal(1 * time.Minute))
}

func TestSyncReconcilerReconcile(t *testing.T) {
	cases := map[string]struct {
		r    reconcile.Reconciler
		eval func(types.Gomega, reconcile.Result, error)
	}{
		"nil reconciler": {
			r: nil,
			eval: func(g types.Gomega, result reconcile.Result, err error) {
				g.Expect(result.IsZero()).To(BeTrue())
				g.Expect(err).NotTo(BeNil())
			},
		},
		"empty reconciler": {
			r: reconcile.Func(emptyReconciler),
			eval: func(g types.Gomega, result reconcile.Result, err error) {
				g.Expect(result.Requeue).To(BeFalse())
				g.Expect(result.RequeueAfter).To(Equal(1 * time.Minute))
				g.Expect(err).To(BeNil())
			},
		},
		"reconciler with error": {
			r: reconcile.Func(reconcilerWithError),
			eval: func(g types.Gomega, result reconcile.Result, err error) {
				g.Expect(result.IsZero()).To(BeTrue())
				g.Expect(err).NotTo(BeNil())
			},
		},
		"reconciler with requeue": {
			r: reconcile.Func(reconcilerWithRequeue),
			eval: func(g types.Gomega, result reconcile.Result, err error) {
				g.Expect(result.Requeue).To(BeTrue())
				g.Expect(result.RequeueAfter).To(Equal(time.Duration(0)))
				g.Expect(err).To(BeNil())
			},
		},
		"reconciler with requeue after": {
			r: reconcile.Func(reconcilerWithRequeueAfter),
			eval: func(g types.Gomega, result reconcile.Result, err error) {
				g.Expect(result.Requeue).To(BeFalse())
				g.Expect(result.RequeueAfter).To(Equal(2 * time.Minute))
				g.Expect(err).To(BeNil())
			},
		},
	}

	for name, item := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			r := NewScheduleReconciler(item.r, func(options *ScheduleOptions) {
				options.SyncPeriod = 1 * time.Minute
			})
			result, err := r.Reconcile(context.Background(), reconcile.Request{})
			item.eval(g, result, err)
		})
	}
}

func emptyReconciler(_ context.Context, _ reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func reconcilerWithError(_ context.Context, _ reconcile.Request) (reconcile.Result, error) {
	result := reconcile.Result{}

	return result, fmt.Errorf("test error")
}

func reconcilerWithRequeue(_ context.Context, _ reconcile.Request) (reconcile.Result, error) {
	result := reconcile.Result{Requeue: true}

	return result, nil
}

func reconcilerWithRequeueAfter(_ context.Context, _ reconcile.Request) (reconcile.Result, error) {
	result := reconcile.Result{RequeueAfter: 2 * time.Minute}

	return result, nil
}
