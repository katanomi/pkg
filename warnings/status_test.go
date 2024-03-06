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

package warnings

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var _ = Describe("Test.GetStatusWarnings.EnsureStatusWarning", func() {
	var (
		status           *duckv1.Status
		key              = "key"
		warning          *WarningRecord
		actual, expected []WarningRecord
	)

	BeforeEach(func() {
		// warning = &WarningRecord{}
		status = &duckv1.Status{}
	})
	JustBeforeEach(func() {
		actual = EnsureStatusWarning(status, key, warning)
		get := GetStatusWarnings(status, key)
		Expect(actual).To(Equal(get))
	})

	When("warning is nil", func() {
		It("should not add the warning", func() {
			Expect(actual).To(HaveLen(0))
		})
	})

	When("warning is not present", func() {
		BeforeEach(func() {
			warning = newWarning("warning 1")
			expected = []WarningRecord{*warning}
		})
		It("should add the warning", func() {
			Expect(actual).To(HaveLen(1))
			Expect(actual).To(Equal(expected))
		})
	})

	When("warning is present", func() {
		BeforeEach(func() {
			warning = newWarning("warning 1")
			expected = []WarningRecord{*warning}
			status.Annotations = map[string]string{
				key: serializeWarnings(expected),
			}
		})
		It("should not add the warning", func() {
			Expect(actual).To(HaveLen(1))
			Expect(actual).To(Equal(expected))
		})
	})

})
