/*
Copyright 2024 The Katanomi Authors.

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

package conformance

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test_testPoint_CheckExternalAssertion", func() {
	var testPoint *testPoint
	var args = []interface{}{1, "1"}

	BeforeEach(func() {
		testPoint = NewTestPoint("test")
	})

	Context("Unregistered custom assertions", func() {
		It("should not panic", func() {
			Expect(func() {
				testPoint.CheckExternalAssertion(args...)
			}).ShouldNot(Panic())
		})
	})

	Context("Registered custom assertions", func() {
		var feature = NewFeatureCase("test-feature")
		BeforeEach(func() {
			testPoint.AddAssertion(func(intVal int, stringVal string) {
				panic("test panic")
			})
			testPoint.bindFeature(feature)
		})

		When("the testcase contains the same labels as the feature", feature.node.Labels(), func() {
			When("parameters are the same as the custom assertion", func() {
				It("should call the custom assertion", func() {
					Expect(func() {
						testPoint.CheckExternalAssertion(args...)
					}).Should(PanicWith("test panic"))
				})
			})
			When("parameters are different from custom assertions", func() {
				It("should panic with reflect error", func() {
					Expect(func() {
						testPoint.CheckExternalAssertion(false)
					}).Should(PanicWith(MatchRegexp("reflect.+")))
				})
			})
		})

		When("the testcase not contains the same labels as the feature", func() {
			It("should not call the custom assertion", func() {
				Expect(func() {
					testPoint.CheckExternalAssertion(args...)
				}).ShouldNot(Panic())
			})
		})
	})
})
