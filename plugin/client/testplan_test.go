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

var _ = Describe("TestPlan", func() {
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

	Describe("List TestPlans", func() {
		It("returns a list of test plans", func() {
			responder := httpmock.NewJsonResponderOrPanic(200, httpmock.File(
				"testdata/fixtures/testplans.json"))

			opt := v1alpha1.TestProjectOptions{
				Project:    "xxx",
				TestPlanID: "123",
			}

			fakeUrl := fmt.Sprintf("%s/projects/%s/testplans", meta.BaseURL, opt.Project)
			url, _ := apis.ParseURL(meta.BaseURL)
			httpmock.RegisterResponder("GET", fakeUrl, responder)
			list, err := pluginClient.TestPlan(meta, secret).List(ctx, &duckv1.Addressable{URL: url}, opt)
			Expect(err).To(BeNil())
			Expect(list.Items).To(HaveLen(1))
		})
	})

	Describe("Get TestPlan", func() {
		It("returns a test plan detail", func() {
			responder := httpmock.NewJsonResponderOrPanic(200, httpmock.File("testdata/fixtures/testplan.json"))

			opt := v1alpha1.TestProjectOptions{
				Project:    "xxx",
				TestPlanID: "387013",
				BuildID:    "456",
			}

			fakeUrl := fmt.Sprintf("%s/projects/%s/testplans/%s?buildID=%s", meta.BaseURL, opt.Project,
				opt.TestPlanID, opt.BuildID)
			url, _ := apis.ParseURL(meta.BaseURL)
			httpmock.RegisterResponder("GET", fakeUrl, responder)
			testPlan, err := pluginClient.TestPlan(meta, secret).Get(ctx, &duckv1.Addressable{URL: url}, opt)
			Expect(err).To(BeNil())
			Expect(testPlan.Spec.ID).To(Equal("387013"))
		})
	})
})
