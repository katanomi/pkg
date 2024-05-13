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

package multicluster

import (
	"context"
	"fmt"

	pkgscheme "github.com/katanomi/pkg/scheme"
	pkgtesting "github.com/katanomi/pkg/testing"
	multiclustertesting "github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/multicluster"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("FilterClusters", func() {
	var (
		clusterRefs         []corev1.ObjectReference
		filteredClusterRefs []corev1.ObjectReference
		ctx                 context.Context

		mgr *ClusterManager
	)

	BeforeEach(func() {
		ctx = context.TODO()
		pkgtesting.MustLoadYaml("testdata/clusterreferences.init.yaml", &clusterRefs)

		mgr = &ClusterManager{
			Concurrent: 1,
			Filters:    []ClusterFilter{},
		}

	})

	JustBeforeEach(func() {
		filteredClusterRefs = mgr.FilterClusters(ctx, clusterRefs)
	})

	Context("empty filters", func() {
		It("should return all clusters", func() {
			Expect(filteredClusterRefs).To(HaveLen(4))
		})
	})

	Context("filter by namespace", func() {
		BeforeEach(func() {
			mgr.Filters = append(mgr.Filters, func(ctx context.Context, objRef corev1.ObjectReference) bool {
				return objRef.Namespace == "default"
			})
		})

		It("should return only cluster in default namespace", func() {
			Expect(filteredClusterRefs).To(HaveLen(2))
		})
	})

})

var _ = Describe("ClusterFilterManager.FilterClusters.CustomResourceDefinitionExists", func() {

	var (
		clusterRefs         []corev1.ObjectReference
		filteredClusterRefs []corev1.ObjectReference
		ctx                 context.Context

		mockCtl       *gomock.Controller
		mockCliGetter *multiclustertesting.MockInterface
		scheme        *runtime.Scheme
		crd           *apiextensionsv1.CustomResourceDefinition

		mgr *ClusterManager
	)

	BeforeEach(func() {
		ctx = context.TODO()
		pkgtesting.MustLoadYaml("testdata/clusterreferences.init.yaml", &clusterRefs)

		mockCtl = gomock.NewController(GinkgoT())
		mockCliGetter = multiclustertesting.NewMockInterface(mockCtl)

		crd = &apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "customresourcedefinitions.apps.k8s.io",
			},
		}

		scheme = runtime.NewScheme()
		apiextensionsv1.AddToScheme(scheme)

		ctx = pkgscheme.WithScheme(ctx, scheme)

		mgr = &ClusterManager{
			Concurrent: 1,
			Filters:    []ClusterFilter{},
		}
	})

	JustBeforeEach(func() {
		mgr.Filters = append(mgr.Filters, CustomResourceDefinitionExists(mockCliGetter, "customresourcedefinitions.apps.k8s.io"))
		filteredClusterRefs = mgr.FilterClusters(ctx, clusterRefs)
	})

	Context("error with client getter", func() {

		BeforeEach(func() {
			mockCliGetter.EXPECT().GetClient(ctx, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("get client err")).AnyTimes()
		})

		It("should return no clusters", func() {
			Expect(filteredClusterRefs).To(HaveLen(0))
		})
	})

	Context("crd exists in cluster", func() {

		BeforeEach(func() {
			mockClientWithCRD := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd).Build()
			mockClientWithoutCRD := fake.NewClientBuilder().WithScheme(scheme).Build()

			gomock.InOrder(
				mockCliGetter.EXPECT().GetClient(ctx, gomock.Any(), gomock.Any()).Return(mockClientWithCRD, nil),
				mockCliGetter.EXPECT().GetClient(ctx, gomock.Any(), gomock.Any()).Return(mockClientWithoutCRD, nil),
				mockCliGetter.EXPECT().GetClient(ctx, gomock.Any(), gomock.Any()).Return(mockClientWithoutCRD, nil),
				mockCliGetter.EXPECT().GetClient(ctx, gomock.Any(), gomock.Any()).Return(mockClientWithCRD, nil),
			)

		})

		It("should return clusters with crd", func() {
			Expect(filteredClusterRefs).To(HaveLen(2))
		})
	})
})
