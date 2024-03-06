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
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// GetStatusWarnings retrieves warnings from a given status annotation.
func GetStatusWarnings(status *duckv1.Status, annotationKey string) []WarningRecord {
	if status == nil || status.Annotations == nil {
		return nil
	}
	return deserializeWarnings(status.Annotations[annotationKey])
}

// EnsureStatusWarning ensures the warning is in the status.
// If the warning already exists, it will be ignored.
func EnsureStatusWarning(status *duckv1.Status, annotationKey string, warning *WarningRecord) []WarningRecord {
	if status == nil || warning == nil {
		return nil
	}
	warnings := GetStatusWarnings(status, annotationKey)
	warnings = AddWarningIfNotPresent(warnings, warning)
	warningStr := serializeWarnings(warnings)
	if len(warningStr) != 0 {
		if status.Annotations == nil {
			status.Annotations = make(map[string]string)
		}
		status.Annotations[annotationKey] = warningStr
	}
	return warnings
}
