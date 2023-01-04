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

package filter

import (
	"context"

	. "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/logging"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

var _ = Describe("Test.GetNamespacesBasedOnFilter", func() {

	var (
		ctx            context.Context
		scheme         *runtime.Scheme
		clt            dynamic.Interface
		err            error
		clusterFilter  ClusterFilter
		actualClusters []corev1.ObjectReference
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		corev1.AddToScheme(scheme)
		ctx = context.TODO()
		clusterFilter = ClusterFilter{}
		objs, err := LoadKubeResourcesAsUnstructured("testdata/clusterfilter.clusters.list.yaml")
		Expect(err).Should(BeNil())
		gvrToListKind := map[schema.GroupVersionResource]string{
			ClusterGVR: "ClusterList",
		}
		clt = dynamicfake.NewSimpleDynamicClientWithCustomListKinds(scheme, gvrToListKind)
		for i, obj := range objs {
			clt.Resource(ClusterGVR).Namespace(obj.GetNamespace()).Create(ctx, &objs[i], metav1.CreateOptions{})
		}
	})

	JustBeforeEach(func() {
		logging.FromContext(context.TODO()).Infow("Debug", "clusterFilter", clusterFilter)
		actualClusters, err = GetClustersBasedOnFilter(ctx, clt, &clusterFilter)
	})

	When("only exist selector", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/clusterfilter.selector.yaml", &clusterFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualClusters).Should(HaveLen(2))
		})
	})

	When("only exist refs", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/clusterfilter.refs.yaml", &clusterFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualClusters).Should(HaveLen(1))
		})
	})

	When("both selector and refs exist and have duplicates", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/clusterfilter.selector.refs.yaml", &clusterFilter)
		})
		It("should filter success", func() {
			Expect(err).Should(BeNil())
			Expect(actualClusters).Should(HaveLen(2))
		})
	})

})
