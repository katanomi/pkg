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

package v1alpha1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"

	knamespace "github.com/katanomi/pkg/namespace"
)

var (
	// ClusterRegistryGroupVersion is the group version for the cluster registry
	ClusterRegistryGroupVersion = schema.GroupVersion{Group: "clusterregistry.k8s.io", Version: "v1alpha1"}
	// ClusterGVR is the group version resource for the cluster registry
	ClusterGVR = ClusterRegistryGroupVersion.WithResource("clusters")

	// passAllClusterFilterRule used to pass all clusters
	passAllClusterFilterRule = ClusterFilterRule{
		Exact: map[string]string{},
	}
)

// GetClustersBasedOnFilter lists clusters based on the clusterFilter criteria
func GetClustersBasedOnFilter(ctx context.Context, clt dynamic.Interface, clusterFilter *ClusterFilter) (clusters []corev1.ObjectReference, err error) {
	if clusterFilter == nil {
		return nil, fmt.Errorf("clusterFilter is nil")
	}
	var ucList *unstructured.UnstructuredList
	// cap of 10 for performance reasons
	clusters = make([]corev1.ObjectReference, 0, 10)
	uClusters := &unstructured.UnstructuredList{}

	defaultNS := clusterFilter.Namespace
	if defaultNS == "" {
		defaultNS = knamespace.NamespaceValue(ctx)
	}
	if defaultNS == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	ucList, err = getClustersByLabelSelector(ctx, clt, clusterFilter.Selector, defaultNS)
	if err != nil {
		return
	}
	uClusters.Items = append(uClusters.Items, ucList.Items...)

	ucList, err = getClustersByRefs(ctx, clt, clusterFilter.Refs, defaultNS)
	if err != nil {
		return
	}
	uClusters.Items = append(uClusters.Items, ucList.Items...)

	// remove duplicates
	uClusters.Items = RemoveDuplicatesFromList(uClusters.Items)

	if clusterFilter.Filter != nil {
		clusters = clusterFilter.Filter.Filter(uClusters.Items)
	} else {
		clusters = passAllClusterFilterRule.Filter(uClusters.Items)
	}

	return
}

func getClustersByLabelSelector(ctx context.Context, clt dynamic.Interface, labelSelector *metav1.LabelSelector, defaultNS string) (uClusters *unstructured.UnstructuredList, err error) {
	uClusters = &unstructured.UnstructuredList{}
	if labelSelector == nil {
		return
	}
	opts := metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(labelSelector),
	}
	clusterList := &unstructured.UnstructuredList{}
	dyclient := clt.Resource(ClusterGVR).Namespace(defaultNS)
	if clusterList, err = dyclient.List(ctx, opts); err != nil {
		return
	}

	for _, cluster := range clusterList.Items {
		if cluster.GetDeletionTimestamp().IsZero() {
			uClusters.Items = append(uClusters.Items, cluster)
		}
	}
	return
}

func getClustersByRefs(ctx context.Context, clt dynamic.Interface, refs []corev1.ObjectReference, defaultNS string) (uClusters *unstructured.UnstructuredList, err error) {
	uClusters = &unstructured.UnstructuredList{}
	if len(refs) == 0 {
		return
	}
	for _, clusterRef := range refs {
		ns := clusterRef.Namespace
		if ns == "" {
			ns = defaultNS
		}
		dyclient := clt.Resource(ClusterGVR).Namespace(ns)
		var cluster *unstructured.Unstructured
		cluster, err = dyclient.Get(ctx, clusterRef.Name, metav1.GetOptions{})
		err = client.IgnoreNotFound(err)
		if err != nil {
			return
		}

		if cluster != nil && cluster.GetDeletionTimestamp().IsZero() {
			uClusters.Items = append(uClusters.Items, *cluster)
		}
	}
	return
}
