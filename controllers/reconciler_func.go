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
	"time"

	"github.com/AlaudaDevops/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// minSyncPeriod defines the minimum synchronization period
const minSyncPeriod = 30 * time.Second

// SyncPeriodResultFunc contruct reconcile.Result to requeue automatically
// the RequeueAfter value is set by configKey in `FeatureFlag`
func SyncPeriodResultFunc(manager *config.Manager, configKey string) ResultFunc {
	return func(_ context.Context, _ reconcile.Request, result *reconcile.Result) error {
		syncPeriod, err := manager.GetFeatureFlag(configKey).AsDuration()
		if err != nil {
			return err
		}

		// do not requeue when sync period is less than or equal to 0
		if syncPeriod.Nanoseconds() <= 0 {
			return nil
		}

		if syncPeriod < minSyncPeriod {
			syncPeriod = minSyncPeriod
		}

		// set the requeue after to the SyncPeriod of the scheduleReconciler
		// ensure that the reconcile run when the result does not require requeue
		if !result.Requeue && result.RequeueAfter <= 0 {
			result.RequeueAfter = syncPeriod
		}
		return nil
	}
}
