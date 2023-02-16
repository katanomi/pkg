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
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/system"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Test_enqueueRequestsMapFunc(t *testing.T) {
	ctx := context.TODO()
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: "ns",
		},
	}

	tests := map[string]struct {
		listFunc ListByFeatureFlagChanged
		obj      client.Object
		want     []reconcile.Request
	}{
		"input nil list func": {
			obj:  &cm,
			want: []reconcile.Request{},
		},
		"match configmap object": {
			obj:      &cm,
			listFunc: func(ctx context.Context) []metav1.Object { return []metav1.Object{&cm} },
			want:     []reconcile.Request{{NamespacedName: types.NamespacedName{Name: "cm", Namespace: "ns"}}},
		},
		"not match configmap object": {
			obj:      &corev1.Secret{},
			listFunc: func(ctx context.Context) []metav1.Object { return []metav1.Object{&cm} },
			want:     []reconcile.Request{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			getFunc := enqueueRequestsConfigMapFunc(ctx, nil, tt.listFunc)
			got := getFunc(tt.obj)

			diff := cmp.Diff(got, tt.want)
			g.Expect(diff).To(BeEmpty())
		})
	}
}

func Test_enqueueRequestsMapFunc_onChange(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.TODO()
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: system.Namespace(),
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	client := fake.NewSimpleClientset()

	watcher := informer.NewInformedWatcher(client, system.Namespace())
	manager := NewManager(watcher, nil, "cm")
	getFunc := enqueueRequestsConfigMapFunc(ctx, manager, func(ctx context.Context) []metav1.Object { return []metav1.Object{&cm} })

	getFunc(&cm)
	config := manager.GetConfig()

	g.Expect(config).NotTo(BeNil(), "config should not nil.")
	g.Expect(cmp.Diff(config.Data, cm.Data)).To(BeEmpty(), "data should be update.")
}

func Test_predicatesUpdateFunc(t *testing.T) {
	oldConfig := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: system.Namespace(),
		},
		Data: map[string]string{
			VersionEnabledFeatureKey: "true",
			PrunerKeepFeatureKey:     "50",
		},
	}

	newConfig := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: system.Namespace(),
		},
		Data: map[string]string{
			VersionEnabledFeatureKey:            "true",
			PrunerKeepFeatureKey:                "50",
			PrunerDelayAfterCompletedFeatureKey: "10s",
		},
	}

	evt := event.UpdateEvent{
		ObjectOld: &oldConfig,
		ObjectNew: &newConfig,
	}

	client := fake.NewSimpleClientset()
	watcher := informer.NewInformedWatcher(client, "cm")
	manager := NewManager(watcher, nil, "cm")

	tests := map[string]struct {
		manager        *Manager
		featureChanged HasFeatureChangedFunc
		evt            event.UpdateEvent
		want           bool
	}{
		"manager is nil": {
			evt:  evt,
			want: false,
		},
		"featureChangedfunc not set": {
			manager: manager,
			evt:     evt,
			want:    false,
		},
		"custom set featureChangedfunc": {
			manager: manager,
			evt:     evt,
			featureChanged: func(new *FeatureFlags, old *FeatureFlags) bool {
				if new == nil || old == nil {
					return false
				}

				newVersionEnable, _ := new.FeatureValue(VersionEnabledFeatureKey).AsBool()
				oldVersionEnable, _ := old.FeatureValue(VersionEnabledFeatureKey).AsBool()
				return newVersionEnable == oldVersionEnable
			},
			want: true,
		},
		"change configmap not set": {
			manager: manager,
			evt:     func() event.UpdateEvent { evt.ObjectNew = nil; return evt }(),
			want:    false,
		},
		"change object is not configmap": {
			manager: manager,
			evt:     func() event.UpdateEvent { evt.ObjectNew = &corev1.Secret{}; return evt }(),
			want:    false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			updateFunc := predicatesUpdateFunc(tt.manager, tt.featureChanged)
			got := updateFunc(tt.evt)

			diff := cmp.Diff(got, tt.want)
			g.Expect(diff).To(BeEmpty())
		})
	}
}
