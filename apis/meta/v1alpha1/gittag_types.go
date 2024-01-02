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

import "errors"

// GitRepositoryTagOption option for repository tag
// Deprecated: use GitTag instead
// +k8s:deepcopy-gen=false
type GitRepositoryTagOption = GitTag

// GitRepositoryTagListOption option for list repository tag
type GitRepositoryTagListOption struct {
	GitRepo
}

// GitTag describe a unique tag
type GitTag struct {
	GitRepo
	// Tag the name of the tag
	Tag string `json:"tag"`
}

// Validate validate the git repo
func (r *GitTag) Validate() error {
	if err := r.GitRepo.Validate(); err != nil {
		return err
	}
	if r.Tag == "" {
		return errors.New("tag is empty")
	}
	return nil
}

type CreateGitTagPayload struct {
	GitTag

	// Ref target of the tag, such as a commit sha/branch/tag
	// sha: c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c
	// branch: refs/heads/branch
	// tag: refs/tags/tag
	// +required
	Ref string `json:"ref"`
	// +optional
	Message string `json:"message"`
}

// Validate validate the git repo
func (r *CreateGitTagPayload) Validate() error {
	if err := r.GitTag.Validate(); err != nil {
		return err
	}
	if r.Ref == "" {
		return errors.New("ref is empty")
	}
	return nil
}
