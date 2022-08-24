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
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetApprovalUsers get the approval users
func (c *Check) GetApprovalUsers() UserApprovals {
	if !c.HasApprover() {
		return nil
	}
	return c.Approval.Users
}

// HasApprover indicates whether the approver exists.
func (c *Check) HasApprover() bool {
	return c != nil && c.Approval.HasApprover()
}

// IsNotSet Indicates that no approver is set.
func (c *Check) IsNotSet() bool {
	return c == nil || c.Approval == nil || len(c.Approval.Users) == 0
}

// ApprovalIsStarted indicates whether the approver is started.
func (cs *ApprovalCheckStatus) ApprovalIsStarted() bool {
	return cs != nil && !cs.ApprovalStartTime.IsZero()
}

// StartWaitingApproval set the approval start time, if it is empty.
func (cs *ApprovalCheckStatus) StartWaitingApproval() {
	if cs == nil {
		return
	}
	if cs.ApprovalStartTime.IsZero() {
		cs.ApprovalStartTime = &metav1.Time{Time: time.Now()}
	}
}

func (users UserApprovals) GetBySubject(subject rbacv1.Subject) *UserApproval {
	if users == nil {
		return nil
	}
	for i := range users {
		user := &users[i]
		if user.Subject == subject {
			return user
		}
	}
	return nil
}

func (as ApprovalStatuses) GetBySubject(subject rbacv1.Subject) *ApprovalStatus {
	if as == nil {
		return nil
	}
	for i := range as {
		user := &as[i]
		if user.Subject == subject {
			return user
		}
	}
	return nil
}
