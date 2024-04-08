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

package auth

import (
	"context"

	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/katanomi/pkg/testing/framework/base"
	"github.com/katanomi/pkg/testing/framework/conformance"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var TestCase = conformance.NewTestCase("Auth")

var (
	TestPointBasicAuth = TestCase.NewTestPoint("BasicAuth")
	TestPointOauth2    = TestCase.NewTestPoint("OAuth2")
)

var CaseSet = TestCase.Build(func(ctx context.Context) {
	caseAuthCheck.DoWithContext(ctx)
})

var caseAuthCheck = P0Case("plugin authentication").
	WithLabels(TestCase).
	WithCondition().WithFunc(func(testContext *TestContext) {

	Context("authentication with oauth2", TestPointOauth2.Labels(testContext.Context), func() {
		statusCode := 200
		It("should be authenticated successful", func() {
			Expect(statusCode).To(Equal(200))
			TestPointOauth2.CheckExternalAssertion(statusCode)
		})
	})

	Context("authentication with basic auth", TestPointBasicAuth.Labels(testContext.Context), func() {
		statusCode := 200
		It("should be authenticated successful", func() {
			Expect(statusCode).To(Equal(200))
			TestPointBasicAuth.CheckExternalAssertion(statusCode)
		})
	})

})
