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
	"testing"
	"time"

	"github.com/katanomi/pkg/config"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSyncPeriodResultFunc(t *testing.T) {
	tests := map[string]struct {
		syncPeriod string
		duration   time.Duration
		err        bool
	}{
		"wrong value": {
			syncPeriod: "abc",
			err:        true,
		},
		"default value": {
			duration: 5 * time.Minute,
		},
		"min value": {
			syncPeriod: "10s",
			duration:   30 * time.Second,
		},
		"normal value": {
			syncPeriod: "1m",
			duration:   1 * time.Minute,
		},
		"zero value": {
			syncPeriod: "0",
			duration:   0 * time.Minute,
		},
	}

	for name, item := range tests {
		t.Run(name, func(t *testing.T) {
			data := map[string]string{}
			if item.syncPeriod != "" {
				data[config.ClusterIntegrationSyncPeriodKey] = item.syncPeriod
			}

			manager := &config.Manager{
				Config: &config.Config{
					Data: data,
				},
			}
			resultFunc := SyncPeriodResultFunc(manager, config.ClusterIntegrationSyncPeriodKey)
			result := &reconcile.Result{}
			err := resultFunc(context.Background(), reconcile.Request{}, result)

			g := NewGomegaWithT(t)
			if item.err {
				g.Expect(err).NotTo(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(result.RequeueAfter).To(Equal(item.duration))
		})
	}
}
