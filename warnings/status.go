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
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// StatusWithWarning is a wrapper around duckv1.Status to provide a way to add warnings to the status.
type StatusWithWarning struct {
	*duckv1.Status
	warningAnnotationKey string
}

// NewStatusWithWarning creates a new StatusWithWarning with the given status and warning annotation key.
func NewStatusWithWarning(status *duckv1.Status, warningAnnotationKey string) *StatusWithWarning {
	if status == nil {
		status = &duckv1.Status{}
	} else {
		// Make a deep copy of the status to avoid modifying the original status.
		status = status.DeepCopy()
	}
	return &StatusWithWarning{
		Status:               status,
		warningAnnotationKey: warningAnnotationKey,
	}
}

// GetStatus retrieves the status.
func (s *StatusWithWarning) GetStatus() *duckv1.Status {
	return s.Status
}

// GetStatusWarnings retrieves warnings from the status.
func (s *StatusWithWarning) GetStatusWarnings() *WarningRecords {
	return NewWarningRecordsFromJSON(s.getRawWarning())
}

// getRawWarning returns the raw warning string from the annotations.
func (s *StatusWithWarning) getRawWarning() string {
	return s.Status.Annotations[s.warningAnnotationKey]
}

// setStatusWarnings sets the warnings to the status.
func (s *StatusWithWarning) setStatusWarnings(warnings *WarningRecords) *StatusWithWarning {
	if warnings == nil {
		delete(s.Status.Annotations, s.warningAnnotationKey)
		return s
	}
	if s.Status.Annotations == nil {
		s.Status.Annotations = make(map[string]string)
	}
	s.Status.Annotations[s.warningAnnotationKey] = warnings.Serialize()
	return s
}

// AddWarning adds a warning to the status.
func (s *StatusWithWarning) AddWarning(warning *WarningRecord) *StatusWithWarning {
	if warning == nil {
		return s
	}
	ws := s.GetStatusWarnings().Add(warning)
	return s.setStatusWarnings(ws)
}

// AddWarningIfNotPresent adds a warning to the status if it is not already present.
func (s *StatusWithWarning) AddWarningIfNotPresent(warning *WarningRecord) *StatusWithWarning {
	if warning == nil {
		return s
	}
	ws := s.GetStatusWarnings().AddIfNotPresent(warning)
	return s.setStatusWarnings(ws)
}

// MakeWarningCondition creates a condition from the warning.
func (s *StatusWithWarning) MakeWarningCondition() *apis.Condition {
	return s.GetStatusWarnings().MakeCondition()
}
