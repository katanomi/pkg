/*
Copyright 2024 The Katanomi Authors.

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

package cluster

import (
	"github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("ExpectKubeObject", func() {
	var (
		clt ctrlclient.Client
	)

	BeforeEach(func() {
		clt = fakeclient.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "config-1",
			},
			Data: map[string]string{
				"a": "a",
				"b": "b",
			},
		}).Build()
	})

	When("clean func is nil", func() {
		It("should comare succeed", func() {
			ExpectKubeObject(&TestContext{Client: clt, Namespace: "default"}, &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "config-1",
				},
			}).Should(testing.DiffEqualTo(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "config-1",
				},
				Data: map[string]string{
					"a": "a",
					"b": "b",
				},
			}))
		})
	})

	When("clean func is not nil", func() {
		It("should comare succeed", func() {
			ExpectKubeObject(&TestContext{Client: clt, Namespace: "default"}, &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "config-1",
				},
			}, testing.KubeObjectDiffClean, func(object interface{}) interface{} {
				delete(object.(*corev1.ConfigMap).Data, "a")
				return object
			}).Should(testing.DiffEqualTo(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "config-1",
				},
				Data: map[string]string{
					"b": "b",
				},
			}))
		})
	})
})
