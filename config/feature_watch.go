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

package config

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ListByFeatureFlagChanged when the function switch is changed, get the object function that triggers reconile.
type ListByFeatureFlagChanged func(ctx context.Context) []metav1.Object

// HasFeatureChangedFunc check whether the function switch of interest has changed.
type HasFeatureChangedFunc func(new *FeatureFlags, old *FeatureFlags) bool

// defaultFeatureChanged the default feature switch comparison function, compares whether all switches have changed.
func defaultFeatureChanged(new *FeatureFlags, old *FeatureFlags) bool {
	return equality.Semantic.DeepEqual(new, old)
}

// WatchFeatureFlagChanged trigger reconcile when the function switch is changed.
func (manager *Manager) WatchFeatureFlagChanged(ctx context.Context, listFunc ListByFeatureFlagChanged, featureChanged HasFeatureChangedFunc) (source.Source, handler.EventHandler, builder.WatchesOption) {

	return &source.Kind{Type: &corev1.ConfigMap{}},
		handler.EnqueueRequestsFromMapFunc(enqueueRequestsConfigMapFunc(ctx, manager, listFunc)),
		// determine whether the function switch has changed, and return true when it changes.
		builder.WithPredicates(predicate.Funcs{
			UpdateFunc: predicatesUpdateFunc(manager, featureChanged),
		})
}

func enqueueRequestsConfigMapFunc(ctx context.Context, manager *Manager, listFunc ListByFeatureFlagChanged) func(client.Object) []reconcile.Request {
	return func(obj client.Object) (reqs []reconcile.Request) {
		reqs = []reconcile.Request{}
		if listFunc == nil {
			return
		}

		configMap, ok := obj.(*corev1.ConfigMap)
		if !ok {
			return
		}

		// Update configuration files manually. the update of the watch may be later than the time trigger.
		if manager != nil && manager.Informer != nil {
			manager.Informer.OnChange(configMap)
		}

		key := types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}
		list := listFunc(ctx)
		for _, t := range list {
			reqs = append(reqs, reconcile.Request{NamespacedName: types.NamespacedName{Name: t.GetName(), Namespace: t.GetNamespace()}})
		}

		log := logging.FromContext(ctx)
		log.Debugw("will enqueue items caused by configmap update", "len(items)", len(list), "configmap", key)
		return reqs
	}
}

func predicatesUpdateFunc(manager *Manager, featureChanged HasFeatureChangedFunc) func(event.UpdateEvent) bool {
	return func(evt event.UpdateEvent) bool {
		var new, old *corev1.ConfigMap
		if evt.ObjectNew != nil {
			// Using the ok mode, when the types do not match, no panic will be triggered.
			new, _ = evt.ObjectNew.(*corev1.ConfigMap)
		}

		if evt.ObjectOld != nil {
			old, _ = evt.ObjectOld.(*corev1.ConfigMap)
		}

		if new == nil || !manager.isSameConfigMap(new) {
			return false
		}

		if featureChanged == nil {
			featureChanged = defaultFeatureChanged
		}

		newConfig := &FeatureFlags{Data: new.Data}
		oldConfig := &FeatureFlags{Data: old.Data}
		return featureChanged(newConfig, oldConfig)
	}
}
