/*
Copyright 2021 The Katanomi Authors.

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

package testing

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

// BuildCondition starts a condition builder object
// useful for unit tests that require validating multiple conditions
func BuildCondition() *ConditionBuilder {
	return &ConditionBuilder{}
}

// ConditionBuilder builds a conditions in a builder pattern
type ConditionBuilder struct {
	apis.Condition
}

// SetType sets the type for the condition
func (c *ConditionBuilder) SetType(ct apis.ConditionType) *ConditionBuilder {
	c.Type = ct
	return c
}

// SetStatus sets the status for the condition
func (c *ConditionBuilder) SetStatus(stat corev1.ConditionStatus) *ConditionBuilder {
	c.Status = stat
	return c
}

// SetReasonMessage sets the message for the condition
func (c *ConditionBuilder) SetReasonMessage(reason, message string, formatKeyValues ...interface{}) *ConditionBuilder {
	c.Reason = reason
	c.Message = fmt.Sprintf(message, formatKeyValues...)
	return c
}

// SetReasonMessage returns a condition
func (c *ConditionBuilder) Done() *apis.Condition {
	return &c.Condition
}
