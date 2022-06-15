/*
Copyright 2022 The Katanomi Authors.

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

package framework

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("GetConfigFromContext", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("provide context without value", func() {
		It("should get nil", func() {
			v := GetConfigFromContext(ctx)
			Expect(v).To(BeNil())
		})
	})

	Context("provide context with a specify value", func() {
		var testValue = "abc"
		BeforeEach(func() {
			ctx = context.WithValue(ctx, configCondition{}, &testValue)
		})

		It("the value should be a string pointer", func() {
			v := GetConfigFromContext(ctx)
			Expect(v).NotTo(BeNil())
			Expect(v).NotTo(Equal(testValue))
		})
	})
})

var _ = Describe("Test NewConfigCondition", func() {
	type configObj struct {
		Name string
	}

	var configName = "test-configmap-name"
	var configContent = "name: abc"
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

	Context("configmap not exist", func() {
		It("should get an error", func() {
			err := NewConfigCondition(configName, &configObj{}).Condition(testCtx)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("configmap exist", func() {
		BeforeEach(func() {
			By("change the namespace witch store the configmap")
			os.Setenv(e2eConfigNSKey, "default")
			DeferCleanup(func() {
				os.Unsetenv(e2eConfigNSKey)
			})

			configmapName := e2eConfigNamePrefix + "-" + configName
			cm := NewTestConfigMap(configmapName, e2eConfigNs, map[string]string{
				"config": configContent,
			})
			err := testCtx.Client.Create(testCtx.Context, cm)
			Expect(err).To(Succeed())
		})

		It("should get the correct config", func() {
			err := NewConfigCondition(configName, &configObj{}).Condition(testCtx)
			Expect(err).To(Succeed())

			v := GetConfigFromContext(testCtx.Context)
			Expect(v).NotTo(BeNil())
			Expect(v).NotTo(Equal(configObj{Name: "abc"}))

			config, err := NewE2EConfig(configName).GetConfig(testCtx.Context, testCtx.Client)
			Expect(err).To(Succeed())
			Expect(config).To(Equal(configContent))
		})
	})
})
