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
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var _ = Describe("Test.StatusWithWarning", func() {
	var (
		status                *duckv1.Status
		statusWithWarning     *StatusWithWarning
		existingWarningRecord *WarningRecord
		warning               *WarningRecord
		warnings              *WarningRecords
		warningAnnotationKey  string
		existingWarningString string
		existingWarningJSON   string
		newWarningJSON        string
	)

	BeforeEach(func() {
		warningAnnotationKey = "warning-annotation-key"
		existingWarningRecord = newWarning("DeprecatedClusterTask", "Existing Warning")
		existingWarningJSON = `{"reason":"DeprecatedClusterTask","message":"Existing Warning"}`
		existingWarningString = "[" + existingWarningJSON + "]"
		//
		warning = newWarning("DeprecatedClusterTask", "New Warning")
		newWarningJSON = `{"reason":"DeprecatedClusterTask","message":"New Warning"}`
		//
		status = &duckv1.Status{
			Annotations: map[string]string{
				warningAnnotationKey: existingWarningString,
			},
		}
		statusWithWarning = NewStatusWithWarning(status, warningAnnotationKey)
	})

	Describe("#GetStatus", func() {
		It("should return the status", func() {
			Expect(statusWithWarning.GetStatus()).To(Equal(status))
		})
	})

	Describe("#GetStatusWarnings", func() {
		Context("when status is not nil", func() {
			It("should return the warning records", func() {
				Expect(statusWithWarning.GetStatusWarnings().Serialize()).To(Equal(existingWarningString))
			})
		})
	})

	Describe("#getRawWarning", func() {
		Context("when status is not nil", func() {
			It("should return the raw warning string", func() {
				Expect(statusWithWarning.getRawWarning()).To(Equal(existingWarningString))
			})
		})
	})

	Describe("#setStatusWarnings", func() {
		var (
			expectedStatus *duckv1.Status
		)

		BeforeEach(func() {
			warnings = NewWarningRecords(existingWarningRecord)
			expectedStatus = &duckv1.Status{
				Annotations: map[string]string{
					warningAnnotationKey: existingWarningString,
				},
			}
		})

		Context("when warnings is nil", func() {
			BeforeEach(func() {
				statusWithWarning.Status = &duckv1.Status{
					Annotations: map[string]string{
						warningAnnotationKey: existingWarningString,
					},
				}
			})

			It("should delete the annotation from the status", func() {
				delete(expectedStatus.Annotations, warningAnnotationKey)
				Expect(statusWithWarning.setStatusWarnings(nil).GetStatus()).To(Equal(expectedStatus))
			})
		})

		Context("when status annotations is nil", func() {
			BeforeEach(func() {
				statusWithWarning.Status = &duckv1.Status{}
				warnings = NewWarningRecords(existingWarningRecord)
				expectedStatus.Annotations = map[string]string{
					warningAnnotationKey: existingWarningString,
				}
			})

			It("should set the warning annotation in the status", func() {
				Expect(statusWithWarning.setStatusWarnings(warnings).GetStatus()).To(Equal(expectedStatus))
			})
		})

		Context("when status annotations is not nil", func() {
			BeforeEach(func() {
				statusWithWarning.Status = &duckv1.Status{
					Annotations: map[string]string{},
				}
			})

			It("should set the warning annotation in the status", func() {
				Expect(statusWithWarning.setStatusWarnings(warnings).GetStatus()).To(Equal(expectedStatus))
			})
		})
	})

	Describe("#AddWarning", func() {
		var (
			expectedStatus   *duckv1.Status
			expectedWarnings *WarningRecords
		)

		BeforeEach(func() {
			expectedStatus = &duckv1.Status{
				Annotations: map[string]string{
					warningAnnotationKey: "[" + existingWarningJSON + "," + newWarningJSON + "]",
				},
			}
			expectedWarnings = NewWarningRecords(existingWarningRecord, warning)
		})

		Context("when warning is nil", func() {
			It("should return itself and not modify anything", func() {
				Expect(statusWithWarning.AddWarning(nil)).To(Equal(statusWithWarning))
				Expect(statusWithWarning.Status).To(Equal(status))
			})
		})

		Context("when status is not nil and warning is not nil", func() {
			It("should add a new warning to the status", func() {
				Expect(statusWithWarning.AddWarning(warning).GetStatus()).To(Equal(expectedStatus))
				Expect(statusWithWarning.GetStatusWarnings()).To(Equal(expectedWarnings))
			})
		})
	})

	Describe("#AddWarningIfNotPresent", func() {
		var (
			expectedStatus   *duckv1.Status
			expectedWarnings *WarningRecords
		)

		BeforeEach(func() {
			expectedStatus = &duckv1.Status{
				Annotations: map[string]string{
					warningAnnotationKey: "[" + existingWarningJSON + "," + newWarningJSON + "]",
				},
			}
			expectedWarnings = NewWarningRecords(existingWarningRecord, warning)
		})

		Context("when warning is nil", func() {
			It("should return itself and not modify anything", func() {
				Expect(statusWithWarning.AddWarningIfNotPresent(nil)).To(Equal(statusWithWarning))
				Expect(statusWithWarning.Status).To(Equal(status))
			})
		})

		Context("when status is not nil and warning is not nil", func() {
			It("should add a new warning to the status if it is not already present", func() {
				Expect(statusWithWarning.AddWarningIfNotPresent(warning).GetStatus()).To(Equal(expectedStatus))
				Expect(statusWithWarning.GetStatusWarnings()).To(Equal(expectedWarnings))
			})
		})
	})

	Describe("#MakeWarningCondition", func() {
		Context("when status is not nil", func() {
			It("should return a warning condition", func() {
				Expect(statusWithWarning.MakeWarningCondition()).To(Equal(&apis.Condition{
					Type:     "Warning",
					Status:   "True",
					Severity: "Warning",
					Reason:   "DeprecatedClusterTask",
					Message:  "Existing Warning",
				}))
			})
		})

	})

})
