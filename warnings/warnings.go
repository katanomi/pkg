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
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
)

// WarningRecords is a list of warning records
type WarningRecords []WarningRecord

// WarningRecord is a record of a warning
type WarningRecord struct {
	// Reason is the reason for the warning
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable message indicating details about the warning.
	// +optional
	Message string `json:"message,omitempty"`

	// Annotations is an unstructured key value map stored with a detail about the warning.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// NewWarningRecords creates a new warning record
func NewWarningRecords(warnings ...*WarningRecord) *WarningRecords {
	w := &WarningRecords{}
	return w.Add(warnings...)
}

// NewWarningRecordsFromJSON creates a new warning record from a json string
func NewWarningRecordsFromJSON(raw string) *WarningRecords {
	w := &WarningRecords{}
	w = w.Deserialize(raw)
	return w
}

// Add adds a warning to the list
func (w *WarningRecords) Add(others ...*WarningRecord) *WarningRecords {
	for _, o := range others {
		w = w.add(o)
	}
	return w
}

func (w *WarningRecords) add(other *WarningRecord) *WarningRecords {
	if other == nil {
		return w
	}
	*w = append(*w, *other)
	return w
}

// AddIfNotPresent adds a warning to the list if it is not already present
func (w *WarningRecords) AddIfNotPresent(others ...*WarningRecord) *WarningRecords {
	for _, o := range others {
		w = w.addIfNotPresent(o)
	}
	return w
}

func (w *WarningRecords) addIfNotPresent(other *WarningRecord) *WarningRecords {
	if other == nil {
		return w
	}
	if w.Has(other) {
		return w
	}
	*w = append(*w, *other)
	return w
}

// Has checks if the warning already exists
func (w *WarningRecords) Has(other *WarningRecord) bool {
	for i := range *w {
		if reflect.DeepEqual(&(*w)[i], other) {
			return true
		}
	}
	return false
}

// Serialize serializes the warnings to a json string
func (w *WarningRecords) Serialize() string {
	raw, err := json.Marshal(w)
	if err != nil {
		logging.FromContext(context.TODO()).Errorw("Failed to serialize warning records", "warnings", w, "error", err)
	}
	return string(raw)
}

// Deserialize deserializes the warnings from the raw json string
func (w *WarningRecords) Deserialize(raw string) *WarningRecords {
	if len(raw) == 0 {
		return w
	}
	*w = []WarningRecord{}
	err := json.Unmarshal([]byte(raw), w)
	if err != nil {
		logging.FromContext(context.TODO()).Errorw("Failed to deserialize warning records", "raw", raw, "error", err)
	}
	return w
}

// MakeCondition generate condition according the warnings
// Example:
//
//	&apis.Condition{
//		Type:     WarningConditionType,
//		Status:   corev1.ConditionTrue,
//		Severity: apis.ConditionSeverityWarning,
//		Reason:   "DeprecatedClusterTask",
//		Message:  "ClusterTask is deprecated\n",
//	}
//
//	&apis.Condition{
//		Type:     WarningConditionType,
//		Status:   corev1.ConditionTrue,
//		Severity: apis.ConditionSeverityWarning,
//		Reason:   "MultipleWarnings",
//		Message:  "1. Warning 1\n2. Warning 2",
//	}
func (w *WarningRecords) MakeCondition() (condition *apis.Condition) {
	if w == nil || len(*w) == 0 {
		return
	}

	condition = &apis.Condition{
		Type:     WarningConditionType,
		Status:   corev1.ConditionTrue,
		Severity: apis.ConditionSeverityWarning,
	}

	if len(*w) == 1 {
		condition.Reason = (*w)[0].Reason
		condition.Message = (*w)[0].Message
		return
	}

	var message strings.Builder
	for i := range *w {
		val := &(*w)[i]
		if i > 0 {
			message.WriteByte('\n')
		}
		fmt.Fprintf(&message, "%d. %s", i+1, val.Message)
	}

	condition.Reason = MultipleWarningsReason
	condition.Message = message.String()
	return
}
