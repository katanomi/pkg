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

type ValidateApprovalFunc func(context.Context, authenticationv1.UserInfo, bool, bool, []*metav1alpha1.ApprovalSpec, []PairOfOldNewCheck) error

// ValidateApproval validates the approval according by the approval spec
// if `allowRepresentOthers` is true, the reqUser can approve on behalf of others
// if `isCreateOperation` is true, the approvalSpec may be nil, skip detection of additional users
func ValidateApproval(ctx context.Context, reqUser authenticationv1.UserInfo, allowRepresentOthers, isCreateOperation bool,
	approvalSpecList []*metav1alpha1.ApprovalSpec, checkList []PairOfOldNewCheck) (err error) {

	defer func() {
		if err != nil {
			log := logging.FromContext(ctx)
			log.Debugw("approval exception detected", "error", err)
		}
	}()

	// check needs to be consistent with the number of spec
	if len(checkList) != len(approvalSpecList) {
		err = fmt.Errorf("internal error #check != #checkSpec")
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
		err = checkApproval(reqUser, allowRepresentOthers, skipAppendCheck, approvalSpec, oldUsers, newUsers)
		if err != nil {
			break
		}
	}
	return
}

func checkApproval(reqUser authenticationv1.UserInfo, allowRepresentOthers, skipAppendCheck bool,
	approvalSpec *metav1alpha1.ApprovalSpec, oldUsers, newUsers metav1alpha1.UserApprovals) (err error) {
	if approvalSpec == nil {
		err = fmt.Errorf("approvalSpec is nil")
		return
	}
	// Cannot add duplicate users
	exists := make(map[rbacv1.Subject]struct{}, len(newUsers))
	for _, user := range newUsers {
		if _, ok := exists[user.Subject]; ok {
			err = fmt.Errorf("approver %q cannot be repeated", user.Subject.Name)
			return
		}
		exists[user.Subject] = struct{}{}
	}

	// Approval users cannot be deleted, and approval results cannot be changed
	for _, oldUser := range oldUsers {
		newUser := newUsers.GetBySubject(oldUser.Subject)
		if newUser == nil {
			err = fmt.Errorf("cannot remove %q from the approval list", oldUser.Subject.Name)
			return
		}
		if oldUser.Input != nil && !reflect.DeepEqual(oldUser.Input, newUser.Input) {
			err = fmt.Errorf("unable to change the approval result for %q from %v to %v", oldUser.Subject.Name, oldUser.Input, newUser.Input)
			return
		}
	}

	// If it is in-order policy, need to be approved in order.
	approvalPolicy := approvalSpec.Policy
	if approvalPolicy == metav1alpha1.ApprovalPolicyInOrder {
		var skippedUser *metav1alpha1.UserApproval
		for i, newUser := range newUsers {
			if skippedUser == nil && newUser.Input == nil {
				skippedUser = &newUsers[i]
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

	// If allow to represent others, no need to check whether the approval behavior is legal.
	if allowRepresentOthers {
		return
	}

	// Prepare a list of legitimate approved users
	exists = make(map[rbacv1.Subject]struct{}, len(approvalSpec.GetApprovalUsers()))
	for _, user := range approvalSpec.GetApprovalUsers() {
		exists[user] = struct{}{}
	}

	for _, newUser := range newUsers {
		// Cannot remove the approver
		oldUser := oldUsers.GetBySubject(newUser.Subject)
		// If it is a create operation, ignore the legality of the new user
		if oldUser == nil && !skipAppendCheck {
			// Only people in the specified list can add the approval result.
			if _, ok := exists[newUser.Subject]; !ok {
				err = fmt.Errorf("%q can not change the approval user list", reqUser.Username)
				return
			}
		}
		// Approval requires verification of identity, cannot approve on behalf of others.
		if ((oldUser == nil || oldUser.Input == nil) && newUser.Input != nil) &&
			!matching.IsRightUser(reqUser, newUser.Subject) {
			err = fmt.Errorf("%q can not approve for user %q", reqUser.Username, newUser.Subject.Name)
			return
		}
	}

	return
}
