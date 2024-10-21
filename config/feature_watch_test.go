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
	"sigs.k8s.io/controller-runtime/pkg/client"
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

			getFunc := enqueueRequestsConfigMapFunc(ctx, tt.listFunc)
			got := getFunc(tt.obj)

			diff := cmp.Diff(got, tt.want)
			g.Expect(diff).To(BeEmpty())
		})
	}
}
