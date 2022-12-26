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

package registry

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	knet "k8s.io/apimachinery/pkg/util/net"

	"github.com/katanomi/pkg/testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	pkgClient "github.com/katanomi/pkg/client"
)

var _ = Describe("Test.RegistrySchemeDetectionBySecret", func() {
	var (
		ctx context.Context
		clt client.Client
		// secret *corev1.Secret
		ref *corev1.ObjectReference
		err error

		registryHost string
		protocols    string
		detect       *RegistrySchemeDetectionBySecret

		mockHttpClient *http.Client
		server         *httptest.Server
	)

	BeforeEach(func() {
		ctx = context.TODO()
		ref = &corev1.ObjectReference{
			Namespace: "default",
			Name:      "secret-name",
		}
		registryHost = ""

		clt = fake.NewClientBuilder().WithScheme(scheme).Build()
		ctx = pkgClient.WithClient(context.Background(), clt)
		detect = NewRegistrySchemeDetectionBySecret(defaultClient, true, true)
	})

	JustBeforeEach(func() {
		protocols, err = detect.DetectScheme(ctx, registryHost)
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	When("resty client is empty", func() {
		BeforeEach(func() {
			detect = NewRegistrySchemeDetectionBySecret(nil, true, true)
		})
		It("should return error", func() {
			Expect(protocols).To(Equal(""))
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).To(BeEquivalentTo("registry client is nil"))

			By("with default scheme")
			protocols = detect.DetectSchemeWithDefault(ctx, registryHost, "https")
			Expect(protocols).To(Equal("https"))
		})
	})

	When("host with scheme https", func() {
		BeforeEach(func() {
			registryHost = "https://127.0.0.1"
		})
		It("should return scheme", func() {
			Expect(protocols).To(Equal("https"))
		})
	})

	When("host with scheme http", func() {
		BeforeEach(func() {
			registryHost = "http://127.0.0.1"
		})
		It("should return scheme", func() {
			Expect(protocols).To(Equal("http"))
		})
	})

	When("secret not found", func() {
		BeforeEach(func() {
			detect = detect.WithClient(clt).WithSecretRef(ref)
		})
		It("should return error", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).To(BeEquivalentTo(`failed to get secret default/secret-name: secrets "secret-name" not found`))
		})
	})

	Context("mock http server", func() {
		BeforeEach(func() {
			rt := knet.SetTransportDefaults(&http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			})
			insecureClient := http.Client{Transport: rt}
			mockHttpClient = &insecureClient

			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}

			server = httptest.NewTLSServer(http.HandlerFunc(handler))
			serverHost := strings.TrimLeft(strings.TrimLeft(server.URL, "http://"), "https://")
			if registryHost == "" {
				registryHost = serverHost
			}

			// inject mock http client
			detect.DefaultRegistrySchemeDetection.httpClient = mockHttpClient
		})
		AfterEach(func() {
			if server != nil {
				server.Close()
			}
		})

		When("secret type is basic", func() {
			BeforeEach(func() {
				detect = detect.WithSecretRef(ref)
				Expect(testing.LoadKubeResources("testdata/secret.basic.yaml", clt)).To(Succeed())
			})
			It("should NOT return error", func() {
				Expect(err).Should(BeNil())
				Expect(protocols).Should(Equal("https"))
			})
		})

		When("secret type is dockerconfig but not matched, no authentication information", func() {
			BeforeEach(func() {
				detect = detect.WithSecretRef(ref)
				Expect(testing.LoadKubeResources("testdata/secret.dockerconfig.yaml", clt)).To(Succeed())
			})
			It("should NOT return error", func() {
				Expect(err).Should(BeNil())
				Expect(protocols).Should(Equal("https"))
			})
		})

		When("secret type is token, no authentication information", func() {
			BeforeEach(func() {
				detect = detect.WithSecretRef(ref)
				Expect(testing.LoadKubeResources("testdata/secret.token.yaml", clt)).To(Succeed())
			})
			It("should NOT return error", func() {
				Expect(err).Should(BeNil())
				Expect(protocols).Should(Equal("https"))
			})
		})
	})
})
