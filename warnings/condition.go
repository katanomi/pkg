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
	"bytes"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

const (
	// WarningConditionType represents a warning condition
	WarningConditionType apis.ConditionType = "Warning"
)

// NewWarningCondition creates a warning condition from the given warnings
func NewWarningCondition(warnings []WarningRecord) (condition *apis.Condition) {
	if len(warnings) == 0 {
		return
	}

	condition = &apis.Condition{
		Type:     WarningConditionType,
		Status:   corev1.ConditionTrue,
		Severity: apis.ConditionSeverityWarning,
	}

	if len(warnings) == 1 {
		condition.Reason = warnings[0].Reason
		condition.Message = warnings[0].Message
		return
	}

	condition.Reason = MultipleWarningsReason
	var message bytes.Buffer
	for i := range warnings {
		message.WriteString(strconv.Itoa(i + 1))
		message.WriteString(". ")
		message.WriteString(warnings[i].Message)
		message.WriteString("\n")
	}
	condition.Message = message.String()
	return
}
