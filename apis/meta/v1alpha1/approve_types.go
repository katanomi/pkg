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

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApprovalSpec Manual approval policy.
type ApprovalSpec struct {
	// Users a list of users that can perform approval for the Stage execution
	Users []rbacv1.Subject `json:"users,omitempty"`

	// Policy defines the policy of approval, such as `Any`, `All`. (Defaults is `Any`)
	// +kubebuilder:default="Any"
	// +kubebuilder:validation:Enum={"Any","All"}
	Policy ApprovalPolicy `json:"policy"`

	// Timeout duration for approval. Once the StageRun starts approval,
	// its execution will wait until it is approved.
	// Once the timeout duration is reached it will automatically fail the StageRun
	Timeout metav1.Duration `json:"timeout"`
}

type UserApprovals []UserApproval

// Approval runtime object for approval process
type Approval struct {
	// Users list of users that can perform approval for the Stage execution
	// have a webhook to ensure that the patch operation will not erase data
	Users UserApprovals `json:"users,omitempty"`
}

// UserApproval specific runtime user record for approval
type UserApproval struct {
	// Approval user data
	rbacv1.Subject `json:",inline"`

	// Operator record the real operator of the approval
	// Exist only when the actual approver is someone else.
	// +optional
	Operator *rbacv1.Subject `json:"operator,omitempty"`

	// Stores the user input for the approval.
	// Once provided may determine the execution of StageRun
	Input *UserApprovalInput `json:"input,omitempty"`
}

// UserApprovalInput user input for a specific approval
type UserApprovalInput struct {
	// Approved user approval decision
	Approved bool `json:"approved,omitempty"`

	// Description for given deicision
	// +optional
	Description string `json:"description,omitempty"`
}

// ApprovalStatus status for approval process
type ApprovalStatus struct {

	// Approval user data
	rbacv1.Subject `json:",inline"`

	// Stores the user input for the approval.
	// Once provided may determine the execution of StageRun
	*UserApprovalInput `json:",inline"`

	// Operator record the real operator of the approval
	// Exist only when the actual approver is someone else.
	// +optional
	Operator *rbacv1.Subject `json:"operator,omitempty"`

	// Time of approval
	ApprovalTime *metav1.Time `json:"approvalTime,omitempty"`
}

// ApprovalPolicy indicate the policy of approval
type ApprovalPolicy string

func (policy ApprovalPolicy) String() string {
	return string(policy)
}

const (
	// ApprovalPolicyAny any one approved it, consider it passed.
	ApprovalPolicyAny ApprovalPolicy = "Any"

	// ApprovalPolicyAll must be approved by all users, consider it passed.
	ApprovalPolicyAll ApprovalPolicy = "All"

	// Not supported at the moment
	// ApprovalPolicyInOrder must be approved by all users in order, consider it passed.
	ApprovalPolicyInOrder ApprovalPolicy = "InOrder"
)
