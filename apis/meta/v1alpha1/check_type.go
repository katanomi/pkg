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

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Check runtime checking process
type Check struct {
	// Approval manual runtime approval data for execution
	Approval *Approval `json:"approval,omitempty"`
}

type ApprovalStatuses []ApprovalStatus

// ApprovalCheckStatus contains approval status
type ApprovalCheckStatus struct {
	// ApprovalStartTime is the time to actually start waiting for approval
	// +optional
	ApprovalStartTime *metav1.Time `json:"approvalStartTime,omitempty"`

	// Approval status for approval process
	// +optional
	Approvals ApprovalStatuses `json:"approvals,omitempty"`
}
