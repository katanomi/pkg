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

package v2

import (
	"context"

	"github.com/jarcoal/httpmock"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var authCheckOption = metav1alpha1.AuthCheckOptions{
	RedirectURL: "test-url",
}

var _ = Describe("Test AuthCheck", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.AuthCheck]()
		httpmock.RegisterResponder(
			"POST",
			"https://example.com/auth/check",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.AuthCheck(context.Background(), authCheckOption)
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})

var _ = Describe("Test AuthToken", func() {
	It("should generate the correct url and got expected response", func() {
		expected := fakeStruct[metav1alpha1.AuthToken]()
		httpmock.RegisterResponder(
			"POST",
			"https://example.com/auth/token",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		got, err := pluginClient.AuthToken(context.Background())
		Expect(err).To(Succeed())
		Expect(diff(got, expected)).To(BeEmpty())
	})
})
