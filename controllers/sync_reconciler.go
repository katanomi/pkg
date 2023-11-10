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

// Importing necessary packages.
import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// minSyncPeriod defines the minimum synchronization period
const minSyncPeriod = 30 * time.Second

type ScheduleOption func(options *ScheduleOptions)

// scheduleReconciler is a wrapper for Reconciler interface from controller-runtime package.
type scheduleReconciler struct {
	reconciler reconcile.Reconciler

	syncPeriod time.Duration
}

type ScheduleOptions struct {
	// define the synchronization interval.
	SyncPeriod time.Duration
}

// defaultOpts set default values for ScheduleOptions
func defaultOpts(opts ScheduleOptions) ScheduleOptions {
	if opts.SyncPeriod <= minSyncPeriod {
		opts.SyncPeriod = minSyncPeriod
	}

	return opts
}

// NewScheduleReconciler creates a new scheduleReconciler with the provided Reconciler and ScheduleOptions.
func NewScheduleReconciler(r reconcile.Reconciler, opts ...ScheduleOption) reconcile.Reconciler {
	options := ScheduleOptions{}
	for _, option := range opts {
		option(&options)
	}
	options = defaultOpts(options)

	return &scheduleReconciler{
		reconciler: r,
		syncPeriod: options.SyncPeriod,
	}
}

// Reconcile is the method that will be called whenever an event occurs that the Reconciler should handle.
// It calls the Reconcile method of the embedded Reconciler and adjusts the RequeueAfter
func (s *scheduleReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	if s.reconciler == nil {
		return reconcile.Result{}, fmt.Errorf("reconciler should not be empty")
	}

	result, err := s.reconciler.Reconcile(ctx, request)
	if err != nil {
		return result, err
	}

	// set the requeue after to the SyncPeriod of the scheduleReconciler
	// ensure that the reconcile run when the result does not require requeue
	if !result.Requeue && result.RequeueAfter <= 0 {
		result.RequeueAfter = s.syncPeriod
	}

	return result, nil
}
