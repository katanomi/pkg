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
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	client2 "github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type Body struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestStoragePluginClientGet(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(200, Body{Message: "Your message", Code: 200})

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewStoragePluginClient(client2.ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects",
		client.Dest(result),
	)

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Your message"))
	g.Expect(result.Code).To(Equal(200))
}

func TestStoragePluginClientPut(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(200, Body{Message: "Changed", Code: 201})

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("PUT", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewStoragePluginClient(client2.ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Put(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects", client.Body(Body{Message: "Changed", Code: 200}), client.Dest(result))

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Changed"))
	g.Expect(result.Code).To(Equal(201))
}

func TestStoragePluginClientError(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(404, nil)

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewStoragePluginClient(client2.ClientOpts(RESTClient))

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects")

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
	client := NewStoragePluginClient(client2.ClientOpts(RESTClient))

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects")

	g.Expect(err).NotTo(BeNil())

	statusError := &errors.StatusError{}
	g.Expect(goerrors.As(err, &statusError)).To(BeTrue())
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}
