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

package approval

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// CheckApprovalResults checks the approval results of the given resources.
// approved is true if passed
// denied is true if rejected
// message is the message of the result
func CheckApprovalResults(ctx context.Context, approvalSpec *metav1alpha1.ApprovalSpec,
	check *metav1alpha1.Check, checkStatus *metav1alpha1.ApprovalCheckStatus) (
	approved, denied bool, message string, newCheckStatus *metav1alpha1.ApprovalCheckStatus) {

	if approvalSpec == nil {
		return false, false, "nil approvalSpec", checkStatus
	}

	// Avoid direct modification of the original value
	newCheckStatus = checkStatus.DeepCopy()
	if newCheckStatus == nil {
		newCheckStatus = &metav1alpha1.ApprovalCheckStatus{}
	}

	log := logging.FromContext(ctx)
	approvalPolicy := approvalSpec.Policy
	checkUsers := check.GetApprovalUsers()
	log.Debugw("check approval result", "policy", approvalPolicy, "approval spec users", approvalSpec.Users,
		"check", check, "approval status", newCheckStatus.Approvals)

	// Get the maximum number of approvers, used to check whether all approvals are passed
	specUserLen := len(approvalSpec.Users)
	checkUserLen := len(check.GetApprovalUsers())
	maxApproverNum := specUserLen
	if checkUserLen > specUserLen {
		maxApproverNum = checkUserLen
	}

	statuses := make([]metav1alpha1.ApprovalStatus, 0, maxApproverNum)
	for _, user := range checkUsers {
		approvalStatus := newCheckStatus.Approvals.GetBySubject(user.Subject)
		// use existing approval results
		if approvalStatus != nil && approvalStatus.UserApprovalInput != nil {
			statuses = append(statuses, *approvalStatus)
			continue
		}
		// still waiting for approval
		if user.Input == nil {
			statuses = append(statuses, metav1alpha1.ApprovalStatus{
				Subject: user.Subject,
			})
			continue
		}
		// record result
		statuses = append(statuses, metav1alpha1.ApprovalStatus{
			Subject:           user.Subject,
			Operator:          user.Operator,
			UserApprovalInput: user.Input,
			ApprovalTime:      &metav1.Time{Time: time.Now()},
		})
	}
	newCheckStatus.Approvals = statuses

	approvedCnt := 0
	approvedMessage := "Approved by"
	for i, user := range checkUsers {
		if user.Input != nil {
			if user.Input.Approved {
				if approvedCnt > 0 {
					approvedMessage += ","
				}
				approvedCnt++
				approvedMessage += fmt.Sprintf(" %q on %s", user.Subject.Name, statuses[i].ApprovalTime.Format(time.RFC3339))
			} else {
				message = fmt.Sprintf("Rejected by %q on %s", user.Subject.Name, statuses[i].ApprovalTime.Format(time.RFC3339))
				denied = true
				break
			}
		}
	}

	if denied {
		return
	}

	if approvalPolicy == metav1alpha1.ApprovalPolicyAny && approvedCnt > 0 {
		approved = true
	}
	if (approvalPolicy == metav1alpha1.ApprovalPolicyAll ||
		approvalPolicy == metav1alpha1.ApprovalPolicyInOrder) &&
		approvedCnt == maxApproverNum {
		approved = true
	}
	if approved {
		message = approvedMessage
		return
	}

	// continue to wait for approval
	return
}
