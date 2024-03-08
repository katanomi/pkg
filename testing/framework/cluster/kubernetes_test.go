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
	"testing"

	"github.com/docker/distribution/context"
	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
			}).Should(ktesting.DiffEqualTo(&corev1.ConfigMap{
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
			}, ktesting.KubeObjectDiffClean, func(object interface{}) interface{} {
				delete(object.(*corev1.ConfigMap).Data, "a")
				return object
			}).Should(ktesting.DiffEqualTo(&corev1.ConfigMap{
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

var _ = Describe("LoadKubeResources", func() {
	var (
		clt      ctrlclient.Client
		file     string
		converts []ktesting.ConvertRuntimeObjctToClientObjectFunc
		err      error
	)

	BeforeEach(func() {
		clt = fakeclient.NewClientBuilder().WithScheme(clientgoscheme.Scheme).Build()
		converts = []ktesting.ConvertRuntimeObjctToClientObjectFunc{}
	})
	JustBeforeEach(func() {
		err = LoadKubeResources(&TestContext{Namespace: "default-e2e", Client: clt}, file, converts...)
	})

	When("file is a kubernetes object that not set namespace", func() {
		BeforeEach(func() {
			file = "./testdata/LoadKubeResources.configmap.yaml"
		})
		It("should create the object in TestContext namespace", func() {
			Expect(err).Should(BeNil())

			configMap := &corev1.ConfigMap{}
			err = clt.Get(context.Background(), ctrlclient.ObjectKey{Namespace: "default-e2e", Name: "default"}, configMap)
			Expect(err).Should(BeNil())
			Expect(configMap.Data["a"]).Should(BeEquivalentTo("1"))
		})
	})

	When("file is a kubernetes object that set namespace", func() {
		BeforeEach(func() {
			file = "./testdata/LoadKubeResources.configmap-withns.yaml"
		})
		It("should create the object in namespace the specified in file", func() {
			Expect(err).Should(BeNil())

			configMap := &corev1.ConfigMap{}
			err = clt.Get(context.Background(), ctrlclient.ObjectKey{Namespace: "default", Name: "default"}, configMap)
			Expect(err).Should(BeNil())
			Expect(configMap.Data["a"]).Should(BeEquivalentTo("1"))
		})
	})

	When("use converts", func() {
		BeforeEach(func() {
			file = "./testdata/LoadKubeResources.configmap.yaml"
			converts = append(converts, func(object runtime.Object) (ctrlclient.Object, error) {
				configmap := object.(*corev1.ConfigMap)
				configmap.Namespace = "default-e2e"
				configmap.Annotations = map[string]string{
					"a": "1",
				}
				return configmap, nil
			})
		})
		It("should create the object according change by convertFunc", func() {
			Expect(err).Should(BeNil())

			configMap := &corev1.ConfigMap{}
			err = clt.Get(context.Background(), ctrlclient.ObjectKey{Namespace: "default-e2e", Name: "default"}, configMap)
			Expect(err).Should(BeNil())
			Expect(configMap.Data["a"]).Should(BeEquivalentTo("1"))
			Expect(configMap.Annotations["a"]).Should(BeEquivalentTo("1"))
		})
	})
})

func TestLoadKubeResources(t *testing.T) {
	clt := fakeclient.NewClientBuilder().WithScheme(clientgoscheme.Scheme).Build()

	g := NewGomegaWithT(t)
	g.Expect(func() {
		MustLoadKubeResources(&TestContext{Namespace: "default", Client: clt}, "./testdata/LoadKubeResources_invalid.yaml")
	}).Should(Panic())
}
