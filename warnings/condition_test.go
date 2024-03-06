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
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

var _ = Describe("Test.NewWarningCondition", func() {
	DescribeTable("should create a warning condition",
		func(warnings []WarningRecord, expectedCondition *apis.Condition) {
			condition := NewWarningCondition(warnings)
			Expect(condition).To(Equal(expectedCondition))
		},
		Entry("should return nil when warnings is empty", []WarningRecord{}, nil),
		Entry("should set reason, message, type, status, and severity when there is only one warning",
			[]WarningRecord{
				{Reason: "reason1", Message: "message1"},
			},
			&apis.Condition{
				Type:     WarningConditionType,
				Status:   corev1.ConditionTrue,
				Severity: apis.ConditionSeverityWarning,
				Reason:   "reason1",
				Message:  "message1",
			},
		),
		Entry("should set reason, message, type, status, and severity when there are multiple warnings",
			[]WarningRecord{
				{Reason: "reason1", Message: "message1"},
				{Reason: "reason2", Message: "message2"},
			},
			&apis.Condition{
				Type:     WarningConditionType,
				Status:   corev1.ConditionTrue,
				Severity: apis.ConditionSeverityWarning,
				Reason:   MultipleWarningsReason,
				Message:  "1. message1\n2. message2\n",
			},
		),
	)
})
