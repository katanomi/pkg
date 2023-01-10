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

	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/logging"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Test.GetNamespacesBasedOnFilter", func() {

	var (
		ctx              context.Context
		scheme           *runtime.Scheme
		clt              client.Client
		err              error
		nsFilter         NamespaceFilter
		actualNamespaces []corev1.Namespace
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		corev1.AddToScheme(scheme)
		ctx = context.TODO()
		nsFilter = NamespaceFilter{}
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		Expect(LoadKubeResources("testdata/namespacefilter.namespaces.list.yaml", clt)).To(Succeed())
	})

	JustBeforeEach(func() {
		logging.FromContext(context.TODO()).Infow("Debug", "nsFilter", nsFilter)
		actualNamespaces, err = GetNamespacesBasedOnFilter(ctx, clt, nsFilter)
	})

	When("only exist selector", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/namespacefilter.selector.yaml", &nsFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualNamespaces).Should(HaveLen(2))
		})
	})

	When("only exist refs", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/namespacefilter.refs.yaml", &nsFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualNamespaces).Should(HaveLen(1))
		})
	})

	When("both selector and refs exist and have duplicates", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/namespacefilter.selector.refs.yaml", &nsFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualNamespaces).Should(HaveLen(2))
		})
	})

})

var _ = Describe("Test.RemoveDuplicatesFromList", func() {

	Context("int list", func() {
		It("should remove duplicates success", func() {
			intList := []int{0, 0, 1, 1, 2}
			actualIntList := RemoveDuplicatesFromList(intList)
			expectedIntList := []int{0, 1, 2}
			Expect(actualIntList).To(Equal(expectedIntList))
		})
	})

	Context("string list", func() {
		It("should remove duplicates success", func() {
			stringList := []string{"", "", "a", "a", "b"}
			actualStringList := RemoveDuplicatesFromList(stringList)
			expectedStringList := []string{"", "a", "b"}
			Expect(actualStringList).To(Equal(expectedStringList))
		})
	})

})
