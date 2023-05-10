/*
Copyright 2023 The Katanomi Authors.

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

// Validate validates GitCreateCommit
func (r *GitCreateCommit) Validate(ctx context.Context) (errs field.ErrorList) {
	return r.Spec.Validate(ctx, field.NewPath("spec"))
}

// Validate validates GitCreateCommitSpec
func (r *GitCreateCommitSpec) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	if r.Message == "" {
		errs = append(errs, field.Invalid(path.Child("message"), "", "commit message is required"))
	}
	if len(r.Actions) == 0 {
		errs = append(errs, field.Invalid(path.Child("actions"), "[]", "actions is required"))
	}

	cnt := 0
	if r.StartBranch != "" {
		cnt++
	}
	if r.StartSHA != "" {
		cnt++
	}
	if r.StartTag != "" {
		cnt++
	}
	if cnt > 1 {
		errs = append(errs, field.Forbidden(path, `only one of startBranch, startSHA OR startTag can be used, not all at the same time.`))
	}
	return
}
