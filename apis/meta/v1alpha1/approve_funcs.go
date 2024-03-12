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

import rbacv1 "k8s.io/api/rbac/v1"

// HasApprover indicates whether the approver exists.
func (a *Approval) HasApprover() bool {
	return a != nil && len(a.Users) != 0
}

// HasApprover indicates whether the approver exists.
func (a *ApprovalSpec) HasApprover() bool {
	return a != nil && len(a.Users) != 0
}

// IsEnabled indicates whether the approval is enabled.
func (a *ApprovalSpec) IsEnabled() bool {
	return a != nil && len(a.Users) != 0
}

// TimeoutEnabled indicates whether the timeout is enabled.
func (a *ApprovalSpec) TimeoutEnabled() bool {
	return a != nil && a.Timeout.Duration > 0
}

// GetApprovalUsers returns the users that can approve.
func (a *ApprovalSpec) GetApprovalUsers() []rbacv1.Subject {
	if a == nil {
		return nil
	}
	return a.Users
}

// IsApproved return input approved status.
func (a *UserApprovalInput) IsApproved() bool {
	if a == nil {
		return false
	}

	if a.Approved != nil && *a.Approved {
		return true
	}
	return false
}

// SetApproved set approved field with param.
func (a *UserApprovalInput) SetApproved(approved bool) {
	a.Approved = &approved
}
