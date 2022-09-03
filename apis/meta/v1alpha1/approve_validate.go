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
	"context"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate validates the Approval policy
func (r ApprovalPolicy) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	switch r {
	case ApprovalPolicyAny, ApprovalPolicyAll, ApprovalPolicyInOrder:
		// valid
	default:
		errs = append(errs, field.NotSupported(path, r, []string{
			string(ApprovalPolicyAny), string(ApprovalPolicyAll), string(ApprovalPolicyInOrder),
		}))
	}
	return
}

// Validate validates the Approval spec
func (r *ApprovalSpec) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	if r == nil {
		return
	}
	errs = append(errs, r.Policy.Validate(ctx, path.Child("policy"))...)
	if r.Timeout.Duration < 0 {
		errs = append(errs, field.Invalid(path.Child("timeout"), r.Timeout.Duration, "should be >= 0"))
	}
	if len(r.Users) == 0 {
		errs = append(errs, field.Invalid(path.Child("users"), r.Users, "expected at least one user"))
	} else {
		path := path.Child("users")
		for idx, user := range r.Users {
			errs = append(errs, ValidateUser(ctx, user, path.Index(idx))...)
		}
	}
	return
}

// ValidateUser validates an user for approval
func ValidateUser(ctx context.Context, user rbacv1.Subject, path *field.Path) (errs field.ErrorList) {
	_ = ctx
	if user.Kind != rbacv1.UserKind && user.Kind != rbacv1.GroupKind {
		errs = append(errs, field.NotSupported(path.Child("kind"), user.Kind, []string{rbacv1.UserKind, rbacv1.GroupKind}))
	}
	if strings.TrimSpace(user.Name) == "" {
		errs = append(errs, field.Required(path.Child("name"), `name is required`))
	}
	if strings.HasPrefix(user.Name, " ") || strings.HasSuffix(user.Name, " ") {
		errs = append(errs, field.Invalid(path.Child("name"), user.Name, "cannot have space prefix or suffix"))
	}
	return
}

// Validate to verify if there are duplicate users
func (r *Approval) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	if r != nil {
		errs = append(errs, r.Users.Validate(ctx, path.Child("users"))...)
	}
	return
}

// Validate to verify if there are duplicate users
func (users UserApprovals) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	set := map[rbacv1.Subject]struct{}{}
	// does not allow repeat the same user in approvals
	for idx, user := range users {
		if _, has := set[user.Subject]; has {
			errs = append(errs, field.Duplicate(path.Index(idx), user.Subject))
		} else {
			set[user.Subject] = struct{}{}
		}
	}
	return
}

// ValidateChange validates changes between fields
func (r *Approval) ValidateChange(ctx context.Context, old *Approval, path *field.Path) (errs field.ErrorList) {
	// already validated inside the ApprovalWebhook
	// here nothing needs to be done
	return
}
