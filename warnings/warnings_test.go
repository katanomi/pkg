/*
Copyright 2024 The AlaudaDevops Authors.

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
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

var (
	w1 = newWarning("reason", "warning-1")
	w2 = newWarning("reason", "warning-2")
	w3 = newWarning("reason", "warning-3")
)

func newWarning(reason, message string) *WarningRecord {
	return &WarningRecord{
		Reason:  reason,
		Message: message,
	}
}

type warningRecordTableEntry struct {
	name           string
	warningRecords WarningRecords
	other          *WarningRecord
	expectedResult WarningRecords
	has            bool
}

var _ = Describe("Test.WarningRecords", func() {

	DescribeTable("Add",
		func(entry warningRecordTableEntry) {
			actualResult := entry.warningRecords.Add(entry.other)
			Expect(actualResult).To(Equal(entry.expectedResult))
			Expect(entry.warningRecords.Has(entry.other)).To(Equal(entry.has))
		},
		Entry("Add when the warning record already exists", warningRecordTableEntry{
			name:           "add a warning record to the list if it already exists",
			warningRecords: NewWarningRecords(w1, w2, w3),
			other:          w3,
			expectedResult: NewWarningRecords(w1, w2, w3, w3),
			has:            true,
		}),
		Entry("Add when the warning record does not exist", warningRecordTableEntry{
			name:           "Add a warning record to the list if it does not exist",
			warningRecords: NewWarningRecords(w1, w2),
			other:          w3,
			expectedResult: NewWarningRecords(w1, w2, w3),
			has:            false,
		}),
	)

	DescribeTable("AddIfNotPresent",
		func(entry warningRecordTableEntry) {
			actualResult := entry.warningRecords.AddIfNotPresent(entry.other)
			Expect(actualResult).To(Equal(entry.expectedResult))
			Expect(entry.warningRecords.Has(entry.other)).To(Equal(entry.has))
		},
		Entry("AddIfNotPresent when the warning record already exists", warningRecordTableEntry{
			name:           "Do not add a warning record to the list if it already exists",
			warningRecords: NewWarningRecords(w1, w2, w3),
			other:          w3,
			expectedResult: NewWarningRecords(w1, w2, w3),
			has:            true,
		}),
		Entry("AddIfNotPresent when the warning record does not exist", warningRecordTableEntry{
			name:           "Add a warning record to the list if it does not exist",
			warningRecords: NewWarningRecords(w1, w2),
			other:          w3,
			expectedResult: NewWarningRecords(w1, w2, w3),
			has:            false,
		}),
	)

	Describe("Serialize", func() {
		It("serializes the warnings to a JSON string", func() {
			warningRecords := NewWarningRecords(&WarningRecord{Reason: "reason", Message: "message"})
			expectedRaw := `[{"reason":"reason","message":"message"}]`
			actualRaw := warningRecords.Serialize()

			Expect(actualRaw).To(Equal(expectedRaw))
		})
	})

	Describe("NewWarningRecordsFromJSON", func() {
		It("deserializes the warnings from a raw JSON string", func() {
			expectedWarningRecords := NewWarningRecords(w1)
			raw := `[{"reason":"reason","message":"warning-1"}]`
			actualWarningRecords := NewWarningRecordsFromJSON(raw)

			Expect(actualWarningRecords).To(Equal(expectedWarningRecords))
		})

		It("returns an empty list if the raw JSON string is empty", func() {
			warningRecords := NewWarningRecords()
			raw := ``
			actualWarningRecords := NewWarningRecordsFromJSON(raw)

			Expect(actualWarningRecords).To(Equal(warningRecords))
		})
	})

	DescribeTable("Has",
		func(entry warningRecordTableEntry) {
			Expect(entry.warningRecords.Has(entry.other)).To(Equal(entry.has))
		},
		Entry("Has when the warning record exists", warningRecordTableEntry{
			name:           "return true if the warning record exists",
			warningRecords: NewWarningRecords(w1, w2, w3),
			other:          w3,
			has:            true,
		}),
		Entry("Has when the warning record does not exist", warningRecordTableEntry{
			name:           "return false if the warning record does not exist",
			warningRecords: NewWarningRecords(w1, w2),
			other:          w3,
			has:            false,
		}),
	)

	DescribeTable("MakeCondition",
		func(records WarningRecords, expectedCondition *apis.Condition) {
			result := records.MakeCondition()
			Expect(result).To(Equal(expectedCondition))
		},
		Entry("empty slice", nil, nil),
		Entry(
			"single warning",
			NewWarningRecords(w1),
			&apis.Condition{
				Type:     WarningConditionType,
				Status:   corev1.ConditionTrue,
				Severity: apis.ConditionSeverityWarning,
				Reason:   "reason",
				Message:  "warning-1",
			},
		),
		Entry(
			"multiple warnings",
			NewWarningRecords(w1, w2, w3),
			&apis.Condition{
				Type:     WarningConditionType,
				Status:   corev1.ConditionTrue,
				Severity: apis.ConditionSeverityWarning,
				Reason:   MultipleWarningsReason,
				Message:  "1. warning-1\n2. warning-2\n3. warning-3",
			},
		),
	)

})
