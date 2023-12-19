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
	"fmt"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func testClassAddress() *duckv1.Addressable {
	u, _ := apis.ParseURL("https://test.com")
	return &duckv1.Addressable{
		URL: u,
	}
}

var _ = Describe("ToolService-CheckAlive", func() {
	var pluginClient *PluginClient
	var meta Meta
	var pluginURL string
	var secret v1.Secret
	var ctx context.Context
	var err error
	var responder httpmock.Responder

	BeforeEach(func() {
		err = nil
		meta = Meta{
			BaseURL: "https://abc.com",
		}
		responder = nil
		secret = secretForTest()
		classAddr := testClassAddress()
		pluginURL = fmt.Sprintf("%s/tools/liveness", classAddr.URL.String())
		pluginClient = NewPluginClient(ClientOpts(defaultClient)).
			WithMeta(meta).WithSecret(secret).WithClassAddress(classAddr)
		ctx = context.Background()
	})

	JustBeforeEach(func() {
		httpmock.RegisterResponder("GET", pluginURL, responder)
		err = pluginClient.NewToolService().CheckAlive(ctx)
	})

	When("response success", func() {
		BeforeEach(func() {
			responder = httpmock.NewJsonResponderOrPanic(200, nil)
		})
		It("", func() {
			Expect(err).Should(Succeed())
		})
	})

	When("response error", func() {
		BeforeEach(func() {
			responder = httpmock.NewJsonResponderOrPanic(500, nil)
		})
		It("returns a test plan detail", func() {
			Expect(err).ShouldNot(Succeed())
		})
	})
})
