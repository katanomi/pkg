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

package multicluster

import (
	"context"
	"sync"

	"github.com/katanomi/pkg/parallel"
	"github.com/katanomi/pkg/scheme"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterFilter for filter cluster reousrces
type ClusterFilter func(ctx context.Context, clusterRef corev1.ObjectReference) bool

// ClusterManager manages clusters collection by filters
type ClusterManager struct {
	Concurrent int
	Filters    []ClusterFilter
}

// FilterClusters returns a filtered list of clusters
func (m *ClusterManager) FilterClusters(ctx context.Context, clusterRefs []corev1.ObjectReference) []corev1.ObjectReference {
	if len(m.Filters) == 0 {
		return clusterRefs
	}

	var filteredClusters []corev1.ObjectReference
	mutex := sync.Mutex{}
	addCluster := func(clusterRef corev1.ObjectReference) {
		mutex.Lock()
		defer mutex.Unlock()
		filteredClusters = append(filteredClusters, clusterRef)
	}

	p := parallel.P(logging.FromContext(ctx), "filterClusters")

	taskFunc := func(objRef corev1.ObjectReference) func() (interface{}, error) {
		return func() (interface{}, error) {
			for _, opt := range m.Filters {
				if !opt(ctx, objRef) {
					return nil, nil
				}
			}
			addCluster(objRef)
			return nil, nil
		}
	}

	for _, objRef := range clusterRefs {
		p.Add(taskFunc(objRef))
	}

	concurrent := m.Concurrent
	if concurrent == 0 {
		concurrent = parallel.DefaultConcurrentNum
	}

	p.Context(ctx).SetConcurrent(concurrent).Do().Wait()

	return filteredClusters
}

// CustomResourceDefinitionExists returns true if the CRD exists in the cluster
func CustomResourceDefinitionExists(cliGetter ClientGetter, CRDName string) ClusterFilter {
	return func(ctx context.Context, clusterRef corev1.ObjectReference) bool {
		logger := logging.FromContext(ctx)

		if cliGetter == nil {
			logger.Errorw("clientGetter is nil")
			return false
		}

		clusterCli, err := cliGetter.GetClient(ctx, &clusterRef, scheme.Scheme(ctx))
		if err != nil {
			logger.Errorw("failed to get cluster client", "cluster", clusterRef, "error", err)
			return false
		}
		crd := &apiextensionsv1.CustomResourceDefinition{}
		err = clusterCli.Get(ctx, client.ObjectKey{Name: CRDName}, crd)
		if err != nil {
			logger.Errorw("failed to get crd", "crd", CRDName, "error", err)
			return false
		}
		logger.Debugw("crd exists in cluster", "crd", CRDName, "cluster", clusterRef.Name)
		return true
	}
}
