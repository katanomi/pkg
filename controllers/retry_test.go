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

package controllers

import (
	"context"

	"github.com/AlaudaDevops/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("CreateOrGetWithRetry", func() {

	var (
		ctx        context.Context
		client     client.Client
		err        error
		object     *v1.Pod
		objectList *v1.PodList
		length     int

		scheme *runtime.Scheme
	)

	BeforeEach(func() {
		ctx = context.TODO()
		scheme = runtime.NewScheme()
		object = &v1.Pod{}
		objectList = &v1.PodList{}
		objectList = &v1.PodList{}
		v1.AddToScheme(scheme)
	})

	JustBeforeEach(func() {
		err = CreateOrGetWithRetry(ctx, client, object)
		Expect(client.List(ctx, objectList)).To(BeNil())
		length = len(objectList.Items)
	})

	When("cluster has no object", func() {

		BeforeEach(func() {
			testing.MustLoadYaml("testdata/pod.yaml", object)
			client = fake.NewClientBuilder().WithScheme(scheme).WithObjects().Build()
		})
		It("err is nil and pod is 1", func() {
			Expect(err).To(BeNil())
			Expect(length).To(Equal(1))
		})
	})

	When("cluster has object exist", func() {

		BeforeEach(func() {
			existObject := &v1.Pod{}
			testing.MustLoadYaml("testdata/pod.yaml", object)
			testing.MustLoadYaml("testdata/pod.yaml", existObject)
			client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(existObject).Build()
		})
		It("err is nil and pod is 1", func() {
			Expect(err).To(BeNil())
			Expect(length).To(Equal(1))
		})
	})
})
