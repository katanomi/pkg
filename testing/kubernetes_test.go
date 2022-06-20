/*
Copyright 2021 The Katanomi Authors.

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

package testing

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("LoadKubeResources", func() {

	var (
		ctx    context.Context
		scheme *runtime.Scheme
		clt    client.Client
		err    error

		configmap    *corev1.ConfigMap
		secret       *corev1.Secret
		configmapKey = client.ObjectKey{
			Namespace: "default",
			Name:      "configmap",
		}
		secretKey = client.ObjectKey{
			Namespace: "default",
			Name:      "secret",
		}
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		corev1.AddToScheme(scheme)
		ctx = context.TODO()
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
	})

	When("load resources without convert", func() {
		BeforeEach(func() {
			Expect(LoadKubeResources("testdata/loadkuberesources.yaml", clt)).To(Succeed())
		})
		It("should get configmap success", func() {
			configmap = new(corev1.ConfigMap)
			err = clt.Get(ctx, configmapKey, configmap)
			Expect(err).Should(BeNil())
		})
		It("should get secret success", func() {
			secret = new(corev1.Secret)
			err = clt.Get(ctx, secretKey, secret)
			Expect(err).Should(BeNil())
		})
	})

	When("load resources with one convert", func() {
		BeforeEach(func() {
			Expect(LoadKubeResources("testdata/loadkuberesources.yaml", clt, convertConfigmap)).To(Succeed())
		})
		It("should list configmap success", func() {
			configmapList := &corev1.ConfigMapList{}
			err = clt.List(ctx, configmapList, client.InNamespace(configmapKey.Namespace))
			Expect(err).Should(BeNil())
			Expect(configmapList.Items).To(HaveLen(1))
		})
		PIt("should list secret failed", func() {
			// This case will failed after upgrading to controller-runtime v0.10.1
			secretList := &corev1.SecretList{}
			err = clt.List(ctx, secretList, client.InNamespace(secretKey.Namespace))
			Expect(err).Should(Not(BeNil()))
			Expect(err.Error()).To(Equal("item[0]: can't assign or convert unstructured.Unstructured into v1.Secret"))
		})
	})

	When("load resources with two converts", func() {
		BeforeEach(func() {
			Expect(LoadKubeResources("testdata/loadkuberesources.yaml", clt, convertConfigmap, convertSecret)).To(Succeed())
		})
		It("should list configmap success", func() {
			configmapList := &corev1.ConfigMapList{}
			err = clt.List(ctx, configmapList, client.InNamespace(configmapKey.Namespace))
			Expect(err).Should(BeNil())
			Expect(configmapList.Items).To(HaveLen(1))
		})
		It("should list secret success", func() {
			secretList := &corev1.SecretList{}
			err = clt.List(ctx, secretList, client.InNamespace(secretKey.Namespace))
			Expect(err).Should(BeNil())
			Expect(secretList.Items).To(HaveLen(1))
		})
	})

})

func convertConfigmap(runtimeObj runtime.Object) (obj client.Object, err error) {
	switch v := runtimeObj.(type) {
	case *corev1.ConfigMap:
		obj = v
	default:
		err = fmt.Errorf("Unsupported gvk: %s", runtimeObj.GetObjectKind().GroupVersionKind())
	}
	return
}

func convertSecret(runtimeObj runtime.Object) (obj client.Object, err error) {
	switch v := runtimeObj.(type) {
	case *corev1.Secret:
		obj = v
	default:
		err = fmt.Errorf("Unsupported gvk: %s", runtimeObj.GetObjectKind().GroupVersionKind())
	}
	return
}
