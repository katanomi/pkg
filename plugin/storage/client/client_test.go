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

package client

import (
	"context"
	goerrors "errors"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	pclient "github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	v1 "knative.dev/pkg/apis/duck/v1"
)

type Body struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var _ = Describe("Test.StoragePluginClient", func() {
	var (
		restClient     *resty.Client
		storageClient  *StoragePluginClient
		opts           []BuildOptions
		err            error
		result         *Body
		fakeUrl        string
		baseUrl        string
		responder      httpmock.Responder
		mockHttpMethod string
	)

	BeforeEach(func() {
		responder, _ = httpmock.NewJsonResponder(200, Body{Message: "Your message", Code: 200})
		restClient = resty.New()
		opts = make([]BuildOptions, 0)
		mockHttpMethod = http.MethodGet
		result = &Body{}
	})

	JustBeforeEach(func() {
		httpmock.Reset()
		httpmock.RegisterResponder(mockHttpMethod, fakeUrl, responder)
		httpmock.ActivateNonDefault(restClient.GetClient())
		url, _ := apis.ParseURL(baseUrl)

		opts = append(opts, WithRestClient(restClient))
		storageClient = NewStoragePluginClient(&v1.Addressable{
			URL: url,
		}, opts...)
	})

	Context("Client HTTP methods", func() {
		BeforeEach(func() {
			opts = []BuildOptions{WithGroupVersion(&schema.GroupVersion{
				Group:   "core",
				Version: "v1alpha1",
			})}
			baseUrl = "https://example.com/api"
			fakeUrl = "https://example.com/api/core/v1alpha1/projects"
		})

		JustAfterEach(func() {
			Expect(err).To(BeNil())
			Expect(result.Message).To(Equal("Your message"))
			Expect(result.Code).To(Equal(200))
		})

		Context("get method", func() {
			BeforeEach(func() {
				mockHttpMethod = http.MethodGet
			})

			It("get correct url and return expected result", func() {
				err = storageClient.Get(context.Background(), "projects",
					pclient.ResultOpts(result),
				)
			})
		})

		Context("post method", func() {
			BeforeEach(func() {
				mockHttpMethod = http.MethodPost
			})

			It("post correct url and return expected result", func() {
				err = storageClient.Post(context.Background(), "projects",
					pclient.ResultOpts(result),
				)
			})
		})

		Context("put method", func() {
			BeforeEach(func() {
				mockHttpMethod = http.MethodPut
			})
			It("put correct url and return expected result", func() {
				err = storageClient.Put(context.Background(), "projects",
					pclient.ResultOpts(result),
				)
			})
		})

		Context("delete method", func() {
			BeforeEach(func() {
				mockHttpMethod = http.MethodDelete
			})
			It("delete correct url and return expected result", func() {
				err = storageClient.Delete(context.Background(), "projects",
					pclient.ResultOpts(result),
				)
			})
		})
	})

	Context("clone client with group version", func() {
		BeforeEach(func() {
			baseUrl = "https://example.com/api"
			fakeUrl = "https://example.com/api/foo/v1/bar"
		})
		JustAfterEach(func() {
			Expect(err).To(BeNil())
			Expect(result.Message).To(Equal("Your message"))
			Expect(result.Code).To(Equal(200))
		})

		It("return new client with group version", func() {
			storageClient = storageClient.ForGroupVersion(&schema.GroupVersion{
				Group:   "foo",
				Version: "v1",
			})
			err = storageClient.Get(context.Background(), "bar", pclient.ResultOpts(result))

		})
	})
})

func TestStoragePluginClientError(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(404, nil)

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	url, _ := apis.ParseURL("https://example.com/api/v1")
	client := NewStoragePluginClient(&v1.Addressable{
		URL: url,
	}, WithRestClient(RESTClient))

	err := client.Get(context.Background(), "projects")

	g.Expect(err).NotTo(BeNil())

	statusError := &errors.StatusError{}
	g.Expect(goerrors.As(err, &statusError)).To(BeTrue())
	g.Expect(statusError.Status().Code).To(Equal(int32(404)))
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

func TestStoragePluginClientErrorReason(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(404, errors.NewNotFound(schema.GroupResource{}, "Project"))

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	url, _ := apis.ParseURL("https://example.com/api/v1")
	client := NewStoragePluginClient(&v1.Addressable{
		URL: url,
	}, WithRestClient(RESTClient))
	err := client.Get(context.Background(), "projects")
	g.Expect(err).NotTo(BeNil())
	statusError := &errors.StatusError{}
	g.Expect(goerrors.As(err, &statusError)).To(BeTrue())
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}
