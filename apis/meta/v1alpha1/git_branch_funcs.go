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
	"encoding/json"
)

// GetBranchStatus return git branch status
func (b *GitBranch) GetBranchStatus() (status *BuildGitBranchStatus) {
	status = &BuildGitBranchStatus{}
	if b == nil {
		b = &GitBranch{}
	}
	status.Name = b.Name
	status.Default = false
	if b.Spec.Default != nil {
		status.Default = *b.Spec.Default
	}
	status.Protected = false
	if b.Spec.Protected != nil {
		status.Protected = *b.Spec.Protected
	}
	status.WebURL = b.Spec.WebURL
	// this was present inside builds and is addaed to keep compatbility
	// but should be removed once the weburl is guaranteed to be in the spec
	if b.Spec.Properties != nil {
		var content map[string]string
		json.Unmarshal(b.Spec.Properties.Raw, &content)
		status.WebURL = content["webURL"]
	}
	return
}
