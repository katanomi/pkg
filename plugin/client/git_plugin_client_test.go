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

package client

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

var _ = Describe("Test.GenerateGitPluginClient", func() {
	var (
		ctx                  context.Context
		secretRef            *corev1.ObjectReference
		gitRepoURL           string
		integrationClassName string
		classAddress         *duckv1.Addressable
		gpclient             *GitPluginClient
		err                  error
	)
	BeforeEach(func() {
		ctx = WithPluginClient(context.TODO(), NewPluginClient())
		secretRef = nil
		gitRepoURL = "https://github.com/katanomi/pkg"
		integrationClassName = "github"
		classAddress = &duckv1.Addressable{}
		classAddress.URL, _ = apis.ParseURL(gitRepoURL)
	})

	JustBeforeEach(func() {
		gpclient, err = GenerateGitPluginClient(ctx, secretRef, gitRepoURL, integrationClassName, classAddress)
	})

	Context("ctx without client and secretRef is not nil", func() {
		BeforeEach(func() {
			secretRef = &corev1.ObjectReference{
				Name: "name",
			}
		})
		It("should return error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Context("invalid git repository url", func() {
		BeforeEach(func() {
			gitRepoURL = "http:// github.com"
		})
		It("should return error", func() {
			Expect(err).NotTo(BeNil())
		})
	})

	Context("valid parameters", func() {
		It("should generate success", func() {
			Expect(err).To(BeNil())
			Expect(gpclient).ToNot(BeNil())
			Expect(gpclient.GetMeta().BaseURL).To(Equal("https://github.com"))
			Expect(gpclient.GitRepo).To(Equal(metav1alpha1.GitRepo{
				Project:    "katanomi",
				Repository: "pkg",
			}))
			Expect(gpclient.ClassAddress).To(Equal(classAddress))
			Expect(gpclient.GetSecret()).To(Equal(corev1.Secret{}))
			Expect(gpclient.IntegrationClassName).To(Equal(integrationClassName))
		})
	})
})
