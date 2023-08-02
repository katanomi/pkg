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
	"fmt"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var _ = Describe("ToolMetadata", func() {
	var (
		pluginClient *PluginClient
		meta         Meta
		secret       v1.Secret
		ctx          context.Context
	)

	BeforeEach(func() {
		pluginClient = NewPluginClient(ClientOpts(defaultClient))
		meta.BaseURL = "https://alauda.io"
		secret = secretForTest()
		ctx = context.Background()
	})

	It("should return tool metadata", func() {
		responder := httpmock.NewJsonResponderOrPanic(200, httpmock.File(
			"testdata/metadata/get_metadata.json"))

		fakeUrl := fmt.Sprintf("%s/tool/metadata", meta.BaseURL)
		url, _ := apis.ParseURL(meta.BaseURL)
		httpmock.RegisterResponder("GET", fakeUrl, responder)

		// Call the GetToolMetadata method.
		toolMeta, err := pluginClient.NewToolMetadata(meta, secret).GetToolMetadata(ctx, &duckv1.Addressable{URL: url})

		// Verify that there is no error.
		Expect(err).ToNot(HaveOccurred())

		// Verify that the tool metadata is not nil.
		Expect(toolMeta).ToNot(BeNil())

		// Verify that the version is "v1.0.0" (the expected tool metadata).
		Expect(toolMeta.Spec.Version).To(Equal("v0.0.1"))
	})

	It("response failed", func() {
		responder := httpmock.NewJsonResponderOrPanic(500, nil)

		fakeUrl := fmt.Sprintf("%s/tool/metadata", meta.BaseURL)
		url, _ := apis.ParseURL(meta.BaseURL)
		httpmock.RegisterResponder("GET", fakeUrl, responder)

		// Call the GetToolMetadata method.
		_, err := pluginClient.NewToolMetadata(meta, secret).GetToolMetadata(ctx, &duckv1.Addressable{URL: url})

		// Verify that there is no error.
		Expect(err).To(HaveOccurred())
	})
})
