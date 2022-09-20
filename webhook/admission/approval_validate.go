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

package admission

import (
	"context"
	"fmt"
	"reflect"

	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"knative.dev/pkg/logging"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/user/matching"
)

type ValidateApprovalFunc func(ctx context.Context, reqUser authenticationv1.UserInfo,
	allowRepresentOthers bool, skipAppendCheck bool, approvalSpecList []*metav1alpha1.ApprovalSpec,
	checkList []PairOfOldNewCheck, triggeredBy *metav1alpha1.TriggeredBy) error

// ValidateApproval validates the approval according by the approval spec
// if `allowRepresentOthers` is true, the reqUser can approve on behalf of others
// if `isCreateOperation` is true, the approvalSpec may be nil, skip detection of additional users
func ValidateApproval(ctx context.Context, reqUser authenticationv1.UserInfo, allowRepresentOthers, isCreateOperation bool,
	approvalSpecList []*metav1alpha1.ApprovalSpec, checkList []PairOfOldNewCheck, triggeredBy *metav1alpha1.TriggeredBy) (err error) {

	log := logging.FromContext(ctx)
	defer func() {
		if err != nil {
			log.Infow("approval exception detected", "error", err)
		}
	}()

	// check needs to be consistent with the number of spec
	if len(checkList) != len(approvalSpecList) {
		err = fmt.Errorf("internal error #check != #checkSpec")
		log.Warnw("validate approval failed", "check", checkList, "checkSpec", approvalSpecList, "error", err)
		return
	}

	for i, checks := range checkList {
		oldUsers := checks[0].GetApprovalUsers()
		newUsers := checks[1].GetApprovalUsers()
		approvalSpec := approvalSpecList[i]
		skipAppendCheck := false
		if approvalSpec == nil && isCreateOperation {
			// If it is a create operation, ignore the legality of the new user
			skipAppendCheck = true
			approvalSpec = &metav1alpha1.ApprovalSpec{}
		}
		c := &checkApproval{ctx, reqUser, allowRepresentOthers, skipAppendCheck, approvalSpec, oldUsers, newUsers, triggeredBy}
		err = c.Check()
		if err != nil {
			break
		}
	}
	return
}

// checkApproval is a helper struct for checking approval
type checkApproval struct {
	ctx                  context.Context
	reqUser              authenticationv1.UserInfo
	allowRepresentOthers bool
	skipAppendCheck      bool

	approvalSpec       *metav1alpha1.ApprovalSpec
	oldUsers, newUsers metav1alpha1.UserApprovals
	triggeredBy        *metav1alpha1.TriggeredBy
}

// Check checks the approval
func (c *checkApproval) Check() (err error) {
	log := logging.FromContext(c.ctx)
	if len(c.oldUsers) == 0 && len(c.newUsers) == 0 {
		log.Debugw("in check approval, no approvalSpec, no approvals, skip checking")
		return nil
	}

	// Cannot add duplicate users
	exists := make(map[rbacv1.Subject]struct{}, len(c.newUsers))
	for _, user := range c.newUsers {
		if _, ok := exists[user.Subject]; ok {
			err = fmt.Errorf("approver %q cannot be repeated", user.Subject.Name)
			return
		}
		exists[user.Subject] = struct{}{}
	}

	// Approval users cannot be deleted, and approval results cannot be changed
	for _, oldUser := range c.oldUsers {
		newUser := c.newUsers.GetBySubject(oldUser.Subject)
		if newUser == nil {
			err = fmt.Errorf("cannot remove %q from the approval list", oldUser.Subject.Name)
			return
		}
		if oldUser.Input != nil && !reflect.DeepEqual(oldUser.Input, newUser.Input) {
			err = fmt.Errorf("unable to change the approval result for %q from %v to %v", oldUser.Subject.Name, oldUser.Input, newUser.Input)
			return
		}
	}

	// If allow to represent others, no need to check whether the approval behavior is legal.
	if c.allowRepresentOthers {
		return
	}

	if c.approvalSpec == nil {
		err = fmt.Errorf("approval spec is nil")
		return
	}

	// If it is in-order policy, need to be approved in order.
	approvalPolicy := c.approvalSpec.Policy
	if approvalPolicy == metav1alpha1.ApprovalPolicyInOrder {
		// general user is not allowed to modify the order
		if orderChanged(c.oldUsers, c.newUsers) {
			err = fmt.Errorf("Approval policy is %q, %q cannot change the order of approvers.", approvalPolicy, c.reqUser.Username)
			return
		}
		var skippedUser *metav1alpha1.UserApproval
		for i, newUser := range c.newUsers {
			if skippedUser == nil && newUser.Input == nil {
				skippedUser = &c.newUsers[i]
				continue
			}
			// After skipping, the approved user is found again
			if skippedUser != nil && newUser.Input != nil {
				err = fmt.Errorf("Approval policy is %q, %q can not approve before %q.",
					approvalPolicy, newUser.Subject.Name, skippedUser.Subject.Name)
				return
			}
		}
	}

	// Prepare a list of legitimate approved users
	exists = make(map[rbacv1.Subject]struct{}, len(c.approvalSpec.GetApprovalUsers()))
	for _, user := range c.approvalSpec.GetApprovalUsers() {
		exists[user] = struct{}{}
	}

	for _, newUser := range c.newUsers {
		// Cannot remove the approver
		oldUser := c.oldUsers.GetBySubject(newUser.Subject)
		// If it is a create operation, ignore the legality of the new user
		if oldUser == nil && !c.skipAppendCheck {
			// Only people in the specified list can add the approval result.
			if _, ok := exists[newUser.Subject]; !ok {
				err = fmt.Errorf("%q can not change the approval user list", c.reqUser.Username)
				return
			}
		}
		// Approval requires verification of identity, cannot approve on behalf of others.
		if (oldUser == nil || oldUser.Input == nil) && newUser.Input != nil {
			if !matching.IsRightUser(c.reqUser, newUser.Subject) {
				err = fmt.Errorf("%q can not approve for user %q", c.reqUser.Username, newUser.Subject.Name)
				return
			}
			// RequiresDifferentApprover if set to true, the user who triggered the StageRun cannot approve, unless an admin
			if c.approvalSpec.RequiresDifferentApprover &&
				c.triggeredBy != nil && c.triggeredBy.User != nil && *c.triggeredBy.User == newUser.Subject {
				err = fmt.Errorf("requiresDifferentApprover is enabled, %q can not approve.", newUser.Subject.Name)
				return
			}
		}
	}

	return
}

// orderChanged returns true if the order of the old and new users is different.
func orderChanged(oldUsers, newUsers metav1alpha1.UserApprovals) bool {
	i, j := 0, 0
	oLen, nLen := len(oldUsers), len(newUsers)
	for i < oLen && j < nLen {
		if oldUsers[i].Subject == newUsers[j].Subject {
			i++
			j++
		} else {
			j++
		}
	}
	return oLen != 0 && i < oLen
}
