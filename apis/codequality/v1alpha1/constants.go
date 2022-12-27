/*
Copyright 2022 The Katanomi Authors.

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

package v1alpha1

import (
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

const (
	// Known code rating metrics

	// Reliability metrics is based on the number of bugs
	// found during code analysis
	CodeRatingTypeReliability = "Reliability"
	// Vunerability metrics is based on the number of security issues
	// found during code analysis
	CodeRatingTypeVunerability = "Vunerability"
	// Maintainability metrics is based on the number of code smells
	// found during code analysis
	CodeRatingTypeMaintainability = "Maintainability"
	// SecurityHotspots metrics is based on the number of possible security issues
	// found during code analysis
	CodeRatingTypeSecurityHotspots = "SecurityHotspots"

	// Succeeded alias for  metav1alpha1.SucceededReason
	Succeeded = metav1alpha1.SucceededReason
	// Failed alias for  metav1alpha1.FailedReason
	Failed = metav1alpha1.FailedReason
)

const (
	// CodeRatingRateA A rating
	CodeRatingRateA = "A"
	// CodeRatingRateB B rating
	CodeRatingRateB = "B"
	// CodeRatingRateC C rating
	CodeRatingRateC = "C"
	// CodeRatingRateD D rating
	CodeRatingRateD = "D"
	// CodeRatingRateE E rating
	CodeRatingRateE = "E"
)
