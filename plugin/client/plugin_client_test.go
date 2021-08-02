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

package client

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type Body struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestPluginClientGet(t *testing.T) {
	// TODO: change this unit test to verify headers sent by client
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(200, Body{Message: "Your message", Code: 200})

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Get(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects",
		client.Dest(result),
		client.Secret(corev1.Secret{Type: corev1.SecretTypeBasicAuth, Data: map[string][]byte{"username": []byte("username")}}),
		client.Meta(Meta{Version: "v1.3.4", BaseURL: "http://plugin.com"}),
	)

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Your message"))
	g.Expect(result.Code).To(Equal(200))
}

func TestPluginClientPut(t *testing.T) {
	g := NewGomegaWithT(t)
	httpmock.Reset()

	responder, _ := httpmock.NewJsonResponder(200, Body{Message: "Changed", Code: 201})

	fakeUrl := "https://example.com/api/v1/projects"
	httpmock.RegisterResponder("PUT", fakeUrl, responder)

	RESTClient := resty.New()
	httpmock.ActivateNonDefault(RESTClient.GetClient())
	client := NewPluginClient(ClientOpts(RESTClient))

	result := &Body{}

	url, _ := apis.ParseURL("https://example.com/api/v1")
	err := client.Put(context.Background(), &duckv1.Addressable{
		URL: url,
	}, "projects", client.Body(Body{Message: "Changed", Code: 200}), client.Dest(result))

	g.Expect(err).To(BeNil())
	g.Expect(result.Message).To(Equal("Changed"))
	g.Expect(result.Code).To(Equal(201))
}
