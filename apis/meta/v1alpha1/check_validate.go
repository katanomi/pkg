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

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate to verify if there are duplicate users
func (r *Check) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	if r != nil {
		errs = append(errs, r.Approval.Validate(ctx, path.Child("approval"))...)
	}
	return
}

// ValidateChange validates changes between fields
func (r *Check) ValidateChange(ctx context.Context, old *Check, path *field.Path) (errs field.ErrorList) {
	if r == nil && old != nil {
		errs = append(errs, field.Forbidden(path, `cannot be unset after creation`))
	}
	if r != nil && old == nil {
		// Allow adding fields, the webhook will be check the permissions.
		// If the check is empty, the controller maybe initialize the check based on the spec information.
	}
	if r != nil && old != nil {
		errs = append(errs, r.Approval.ValidateChange(ctx, old.Approval, path.Child("approval"))...)
	}
	return
}
