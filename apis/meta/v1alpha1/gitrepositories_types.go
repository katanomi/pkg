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
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	GitRepoGVK     = GroupVersion.WithKind("GitRepository")
	GitRepoListGVK = GroupVersion.WithKind("GitRepositoryList")
)

// GitRepository repository
type GitRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GitRepositorySpec `json:"spec"`
}

// GitRepositorySpec spec for repository
type GitRepositorySpec struct {
	// Name repo name
	Name string `json:"name"`
	// HtmlURL repo URL
	HtmlURL string `json:"htmlUrl"`
	// HttpCloneURL clone with http
	HttpCloneURL string `json:"httpCloneUrl"`
	// SshCloneURL clone with ssh
	SshCloneURL string `json:"sshCloneUrl"`
	// DefaultBranch main branch name
	DefaultBranch string `json:"defaultBranch"`
	// CreatedAt repo create time
	CreatedAt metav1.Time `json:"createdAt"`
	// UpdatedAt repo last update time
	UpdatedAt metav1.Time `json:"updatedAt"`
	// Owner repo owner
	Owner      GitUserBaseInfo       `json:"owner"`
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

// GitRepositoryList list of repo
type GitRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []GitRepository `json:"items"`
}

// GitRepositoryResourceAttributes returns a ResourceAttribute object to be used in a filter
func GitRepositoryResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		// TODO: modify accordingly
		Resource: "gitrepositories",
		Verb:     verb,
	}
}
