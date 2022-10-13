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

package secret

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	pkgClient "github.com/katanomi/pkg/client"
	kerrors "github.com/katanomi/pkg/errors"
	pkgnamespace "github.com/katanomi/pkg/namespace"
	"github.com/katanomi/pkg/testing"
)

var _ = Describe("Test.GetSecretByRefOrLabel", func() {
	var (
		ctx    context.Context
		clt    client.Client
		secret *corev1.Secret
		ref    *corev1.ObjectReference
		err    error
	)

	BeforeEach(func() {
		ctx = context.TODO()
		ref = &corev1.ObjectReference{
			Namespace: "default",
			Name:      "secret-name",
		}
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = pkgClient.WithClient(context.Background(), clt)
	})

	JustBeforeEach(func() {
		secret, err = GetSecretByRefOrLabel(ctx, clt, ref)
	})

	When("clt is empty", func() {
		BeforeEach(func() {
			clt = nil
		})
		It("should return error", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(Equal(kerrors.ErrNilPointer))
			Expect(secret).Should(BeNil())
		})
	})

	Context("get secret by ref", func() {
		BeforeEach(func() {
			Expect(testing.LoadKubeResources("testdata/secret.no.labels.yaml", clt)).To(Succeed())
		})
		When("ref namespace is not empty", func() {
			It("should return secret", func() {
				Expect(err).Should(BeNil())
				Expect(secret.GetName()).Should(Equal("secret-name"))
			})
		})
		When("ref namespace is empty", func() {
			BeforeEach(func() {
				ref.Namespace = ""
				ctx = pkgnamespace.WithNamespace(ctx, "default")
			})
			It("should return secret", func() {
				Expect(err).Should(BeNil())
				Expect(secret.GetName()).Should(Equal("secret-name"))
			})
		})
		When("ref and ctx namespace both empty", func() {
			BeforeEach(func() {
				ref.Namespace = ""
				ctx = pkgnamespace.WithNamespace(ctx, "")
			})
			It("should return error", func() {
				Expect(err).ShouldNot(BeNil())
				Expect(err).Should(Equal(namespaceIsEmpty))
			})
		})
	})

	Context("get secret by labels", func() {
		BeforeEach(func() {
			Expect(testing.LoadKubeResources("testdata/secret.has.sync.labels.yaml", clt)).To(Succeed())
		})
		When("there is only one secret but has sync label", func() {
			It("should not found any secret", func() {
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("there are two secret and one of them meets the requirements", func() {
			BeforeEach(func() {
				Expect(testing.LoadKubeResources("testdata/secret.has.labels.yaml", clt)).To(Succeed())
			})
			It("should return secret", func() {
				Expect(err).Should(BeNil())
				Expect(secret.GetName()).Should(Equal("secret-has-labels"))
			})
		})
	})
})
