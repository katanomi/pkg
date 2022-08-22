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
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var _ = Describe("TestCase", func() {
	var pluginClient *PluginClient
	var meta Meta
	var secret v1.Secret
	var ctx context.Context

	BeforeEach(func() {
		pluginClient = NewPluginClient(ClientOpts(defaultClient))
		meta.BaseURL = "https://alauda.io"
		secret = secretForTest()
		ctx = context.Background()
	})

	Describe("List TestCases", func() {
		It("returns a list of test cases", func() {
			responder := httpmock.NewJsonResponderOrPanic(200, httpmock.File(
				"testdata/fixtures/testcases.json"))

			opt := v1alpha1.TestProjectOptions{
				Project:    "xxx",
				TestPlanID: "123",
				BuildID:    "456",
			}

			fakeUrl := fmt.Sprintf("%s/projects/%s/testplans/%s/testcases?buildID=%s", meta.BaseURL, opt.Project,
				opt.TestPlanID, opt.BuildID)
			url, _ := apis.ParseURL(meta.BaseURL)
			httpmock.RegisterResponder("GET", fakeUrl, responder)
			list, err := pluginClient.TestCase(meta, secret).List(ctx, &duckv1.Addressable{URL: url}, opt)
			Expect(err).To(BeNil())
			Expect(list.Items).To(HaveLen(1))
		})
	})

	Describe("Get TestCase", func() {
		It("returns a test case detail", func() {
			responder := httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/fixtures/testcase.json"))

			opt := v1alpha1.TestProjectOptions{
				Project:    "xxx",
				TestPlanID: "123",
				TestCaseID: "ACP-82296",
				BuildID:    "456",
			}

			fakeUrl := fmt.Sprintf("%s/projects/%s/testplans/%s/testcases/%s?buildID=%s", meta.BaseURL, opt.Project,
				opt.TestPlanID, opt.TestCaseID, opt.BuildID)
			url, _ := apis.ParseURL(meta.BaseURL)
			httpmock.RegisterResponder("GET", fakeUrl, responder)
			testCase, err := pluginClient.TestCase(meta, secret).Get(ctx, &duckv1.Addressable{URL: url}, opt)
			Expect(err).To(BeNil())
			Expect(testCase.Spec.ID).To(Equal("ACP-82296"))
		})
	})
})
