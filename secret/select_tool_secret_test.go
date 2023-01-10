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

package secret

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/logging"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kclient "github.com/katanomi/pkg/client"
	. "github.com/katanomi/pkg/testing"
)

var _ = Describe("Test.SelectToolSecretByRefOrLabelOrURL", func() {

	var (
		ctx       context.Context
		clt       client.Client
		ns, url   string
		secret    *corev1.Secret
		secretRef *corev1.ObjectReference
		err       error
	)

	BeforeEach(func() {
		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = kclient.WithClient(context.TODO(), clt)
		ctx = logging.WithLogger(ctx, logger)
		ns = "default"
		url = "https://github.com/katanomi/spec"
		secret = &corev1.Secret{}
		secretRef = nil
	})

	JustBeforeEach(func() {
		secret, err = SelectToolSecretByRefOrLabelOrURL(ctx, ns, url, secretRef)
	})

	Context("secret ref is specified", func() {
		BeforeEach(func() {
			secretRef = &corev1.ObjectReference{
				Name: "secret-name",
			}
		})
		When("secret not matched", func() {
			It("should return error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`secret default/secret-name not exist or matching by labels 'map[core.kubernetes.io/namespace:default core.kubernetes.io/secret:secret-name]'`))
			})
		})
		When("secret matched", func() {
			BeforeEach(func() {
				MustLoadYaml("testdata/select.git.secret.yaml", secret)
				Expect(clt.Create(ctx, secret)).To(Succeed())
			})
			It("should NOT return error", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(secret).ShouldNot(BeNil())
				Expect(secret.Name).To(Equal("secret-name"))
			})
		})
	})

	Context("secret ref is not specified", func() {
		When("secret not matched", func() {
			It("should return error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(noSecretSelected))
			})
		})
		When("secret matched", func() {
			BeforeEach(func() {
				MustLoadYaml("testdata/select.git.secret.yaml", secret)
				Expect(clt.Create(ctx, secret)).To(Succeed())
			})
			It("should NOT return error", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(secret).ShouldNot(BeNil())
				Expect(secret.Name).To(Equal("secret-name"))
			})
		})
	})

	Context("secret ref is not specified and url is empty", func() {
		BeforeEach(func() {
			url = ""
		})
		It("should return error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(noSecretSelected))
		})
	})

})
