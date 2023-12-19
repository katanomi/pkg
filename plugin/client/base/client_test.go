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

package base

import (
	"context"
	goerrors "errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type Body struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var testHeaderData = map[string]string{
	PluginAuthHeader:   "kubernetes.io/basic-auth",
	PluginSecretHeader: "eyJ1c2VybmFtZSI6ImRYTmxjbTVoYldVPSJ9",
	PluginMetaHeader:   "eyJ2ZXJzaW9uIjoidjEuMy40IiwiYmFzZVVSTCI6Imh0dHA6Ly9wbHVnaW4uY29tIn0=",
	"Content-Type":     "application/json",
}

func validateRequestHeader(g Gomega, r *http.Request, wantHeaders []string) {
	for _, headerKey := range wantHeaders {
		g.Expect(r.Header.Get(headerKey)).To(Equal(testHeaderData[headerKey]))
	}
}
func TestPluginClientGet(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, func(r *http.Request) (*http.Response, error) {
		validateRequestHeader(g, r, []string{PluginAuthHeader, PluginSecretHeader, PluginMetaHeader, "Content-Type"})
		return httpmock.NewJsonResponse(200, Body{Message: "Your message", Code: 200})
	})

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects",
		ResultOpts(result),
		SecretOpts(corev1.Secret{Type: corev1.SecretTypeBasicAuth, Data: map[string][]byte{"username": []byte("username")}}),
		MetaOpts(Meta{Version: "v1.3.4", BaseURL: "http://plugin.com"}),
	)

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Your message"))
	g.Expect(result.Code).To(Equal(200))
}

func TestPluginClientPut(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("PUT", fakeUrl, func(r *http.Request) (*http.Response, error) {
		validateRequestHeader(g, r, []string{"Content-Type"})
		g.Expect(io.ReadAll(r.Body)).To(Equal([]byte(`{"message":"Changed","code":200}`)))
		return httpmock.NewJsonResponse(200, Body{Message: "Changed", Code: 201})
	})

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Put(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects", BodyOpts(Body{Message: "Changed", Code: 200}), ResultOpts(result))

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Changed"))
	g.Expect(result.Code).To(Equal(201))
}

func TestPluginClientError(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(404, nil)

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

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

func TestPluginClientErrorReason(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(404, errors.NewNotFound(schema.GroupResource{}, "Project"))

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects")

	g.Expect(err).NotTo(BeNil())

	statusError := &errors.StatusError{}
	g.Expect(goerrors.As(err, &statusError)).To(BeTrue())
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

func TestPluginClientClone(t *testing.T) {
	g := NewGomegaWithT(t)

	original := &PluginClient{}
	clone := original.Clone()

	g.Expect(clone).To(Equal(original))
	g.Expect(fmt.Sprintf("%p", original)).NotTo(Equal(fmt.Sprintf("%p", clone)))
}
