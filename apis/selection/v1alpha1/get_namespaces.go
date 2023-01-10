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

	"github.com/katanomi/pkg/hash"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetNamespacesBasedOnFilter lists namespaces based on the namespaceFilter criteria
func GetNamespacesBasedOnFilter(ctx context.Context, clt client.Client, nsFilter NamespaceFilter) (namespaces []corev1.Namespace, err error) {
	var nsList []corev1.Namespace
	// cap of 10 for performance reasons
	namespaces = make([]corev1.Namespace, 0, 10)

	nsList, err = getNamespacesByLabelSelector(ctx, clt, nsFilter.Selector)
	if err != nil {
		return nil, err
	}
	namespaces = append(namespaces, nsList...)

	nsList, err = getNamespacesByRefs(ctx, clt, nsFilter.Refs)
	if err != nil {
		return nil, err
	}
	namespaces = append(namespaces, nsList...)

	// remove duplicates
	namespaces = RemoveDuplicatesFromList(namespaces)

	if nsFilter.Filter != nil {
		namespaces = nsFilter.Filter.Filter(namespaces)
	}

	return
}

func getNamespacesByLabelSelector(ctx context.Context, clt client.Client, labelSelector *metav1.LabelSelector) (namespaces []corev1.Namespace, err error) {
	if labelSelector == nil {
		return
	}

	var selector labels.Selector
	selector, err = metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return
	}
	nsList := &corev1.NamespaceList{}
	if err = clt.List(ctx, nsList, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return
	}

	for _, ns := range nsList.Items {
		if ns.DeletionTimestamp.IsZero() {
			namespaces = append(namespaces, ns)
		}
	}
	return
}

func getNamespacesByRefs(ctx context.Context, clt client.Client, refs []corev1.ObjectReference) (namespaces []corev1.Namespace, err error) {
	if len(refs) == 0 {
		return
	}
	for _, nsref := range refs {
		ns := &corev1.Namespace{}
		ref := nsref
		err = clt.Get(ctx, client.ObjectKey{Name: ref.Name}, ns)
		err = client.IgnoreNotFound(err)
		if err != nil {
			return
		}

		if ns.Name != "" && ns.DeletionTimestamp.IsZero() {
			namespaces = append(namespaces, *ns)
		}
	}
	return
}

// RemoveDuplicatesFromList removes duplicate items from the list
func RemoveDuplicatesFromList[T any](anyList []T) []T {
	rets := make([]T, 0, len(anyList))
	exist := make(map[string]struct{}, len(anyList))
	for _, item := range anyList {
		itemHash := hash.ComputeHash(item)
		if _, ok := exist[itemHash]; !ok {
			exist[itemHash] = struct{}{}
			rets = append(rets, item)
		}
	}
	return rets
}
