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
)

func newWarning(message string) *WarningRecord {
	return &WarningRecord{
		Message: message,
	}
}

var _ = Describe("Test.AddWarningIfNotPresent", func() {
	var (
		tips1 = *newWarning("warning 1")
		tips2 = *newWarning("warning 2")
		tips3 = *newWarning("warning 3")
	)

	DescribeTable("should add a warning if not present",
		func(warnings []WarningRecord, add *WarningRecord, expected []WarningRecord) {
			actual := AddWarningIfNotPresent(warnings, add)
			Expect(actual).To(Equal(expected))
		},
		Entry("when the warning is nil",
			[]WarningRecord{tips1}, nil, []WarningRecord{tips1}),
		Entry("when the warning is not present",
			[]WarningRecord{tips1}, &tips2, []WarningRecord{tips1, tips2}),
		Entry("when the warning is present",
			[]WarningRecord{tips1, tips2, tips3}, &tips3, []WarningRecord{tips1, tips2, tips3}),
	)
})
