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

// CreateBranchParams params for create file in server
type CreateRepoFileParams struct {
	// Branch target branch to create file
	Branch string `json:"branch"`
	// Message commit message for create the file
	Message string `json:"message"`
	// Content must be base64 encoded
	Content []byte `json:"content"`
}

// CreateRepoFilePayload option for create file and commit + push
type CreateRepoFilePayload struct {
	GitRepo
	CreateRepoFileParams
	FilePath string `json:"filepath"`
}

// CreateBranchParams params for create branch in server
type CreateBranchParams struct {
	// Branch new branch name
	Branch string `json:"branch"`
	// Ref source branch name
	Ref string `json:"ref"`
}

// CreateBranchPayload payload for create branch
type CreateBranchPayload struct {
	GitRepo
	CreateBranchParams
}

// CreatePullRequestPayload option for create PullRequest
type CreatePullRequestPayload struct {
	Source      GitBranchBaseInfo `json:"source"`
	Target      GitBranchBaseInfo `json:"target"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
}

// CreatePullRequestCommentPayload payload for create pr's Comment
type CreatePullRequestCommentPayload struct {
	GitRepo
	CreatePullRequestCommentParam
	Index int `json:"index"`
}

// CreatePullRequestCommentParam param for create pr's comment
type CreatePullRequestCommentParam struct {
	Body string `json:"body"`
}
