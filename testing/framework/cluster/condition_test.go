/*
Copyright 2022 The AlaudaDevops Authors.

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
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Test ConditionFunc", func() {
	It("should get the value set by the ConditionFun", func() {
		testCtx := &TestContext{}
		err := ConditionFunc(func(testCtx *TestContext) error {
			testCtx.Namespace = "xxx"
			return nil
		}).Condition(testCtx)
		Expect(err).To(Succeed())
		Expect(testCtx.Namespace).To(Equal("xxx"))
	})
})

var _ = Describe("TestCondition", func() {
	var testCtx *TestContext

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		err := corev1.AddToScheme(scheme)
		Expect(err).To(Succeed())

		testCtx = &TestContext{}
		testCtx.Context = context.TODO()
		testCtx.Namespace = "ccc"
		testCtx.Client = fake.NewClientBuilder().WithScheme(scheme).Build()
	})

	It("should get the test namespace", func() {
		condition := &TestNamespaceCondition{}
		err := condition.Condition(testCtx)
		Expect(err).To(Succeed())

		ns := &corev1.Namespace{}
		key := types.NamespacedName{Name: testCtx.Namespace}
		err = testCtx.Client.Get(testCtx.Context, key, ns)
		Expect(err).To(Succeed())
		Expect(ns.Name).To(Equal(testCtx.Namespace))
	})

	Context("create a configmap", func() {
		var cm *corev1.ConfigMap
		BeforeEach(func() {
			cm = NewTestConfigMap("aa", "default", nil)
			err := testCtx.Client.Create(testCtx.Context, cm)
			Expect(err).To(Succeed())
		})

		When("rollback action", func() {
			It("should roll back successfully", func() {
				key := types.NamespacedName{
					Namespace: cm.Namespace,
					Name:      cm.Name,
				}
				newCm := &corev1.ConfigMap{}
				err := testCtx.Client.Get(testCtx.Context, key, newCm)
				Expect(err).To(Succeed())

				MustRollback(testCtx, cm)

				newCm2 := &corev1.ConfigMap{}
				err = testCtx.Client.Get(testCtx.Context, key, newCm2)
				Expect(errors.IsNotFound(err)).To(BeTrue())
			})
		})
	})

	Context("rollback but resource not found", func() {
		cm := NewTestConfigMap("aa", "default", nil)

		When("rollback action", func() {
			It("should roll back successfully", func() {
				key := types.NamespacedName{
					Namespace: cm.Namespace,
					Name:      cm.Name,
				}
				newCm := &corev1.ConfigMap{}
				err := testCtx.Client.Get(testCtx.Context, key, newCm)
				Expect(errors.IsNotFound(err)).To(BeTrue())

				MustRollback(testCtx, cm)
			})
		})
	})

})
