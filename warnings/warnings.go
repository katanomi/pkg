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
	"encoding/json"
	"reflect"
)

// WarningRecord is a record of a warning
type WarningRecord struct {
	// The reason for the warning tips
	// +optional
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the warning.
	// +optional
	Message string `json:"message,omitempty"`

	// Annotations is an unstructured key value map stored with a detail about the warning.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AddWarningIfNotPresent adds the warning to the list if it is not present
func AddWarningIfNotPresent(warnings []WarningRecord, warning *WarningRecord) []WarningRecord {
	if warning == nil || hasWarning(warnings, warning) {
		return warnings
	}
	warnings = append(warnings, *warning)
	return warnings
}

// serializeWarnings serializes the warnings to a json string
func serializeWarnings(warnings []WarningRecord) string {
	raw, _ := json.Marshal(warnings)
	return string(raw)
}

// deserializeWarnings deserializes the warnings from the raw json string
func deserializeWarnings(raw string) (warnings []WarningRecord) {
	if raw == "" {
		return
	}
	warnings = []WarningRecord{}
	json.Unmarshal([]byte(raw), &warnings)
	return
}

// hasWarning checks if the warning already exists
func hasWarning(warnings []WarningRecord, warning *WarningRecord) bool {
	for i := range warnings {
		if reflect.DeepEqual(&warnings[i], warning) {
			return true
		}
	}
	return false
}
