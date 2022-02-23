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
)

var UserGVK = GroupVersion.WithKind("User")
var UserListGVK = GroupVersion.WithKind("UserList")

// User object for plugin
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IssueSpec `json:"spec"`
}

// UserSpec for Issue
type UserSpec struct {
	// user id
	Id string `json:"id,omitempty"`

	// user name
	Name string `json:"name,omitempty"`

	//user email
	Email string `json:"email,omitempty"`

	// add more field...
}

// UserList list of user
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	ListMeta        `json:"metadata,omitempty"`

	Items []User `json:"items"`
}

// UserResourceAttributes returns a ResourceAttribute object to be used in a filter
func UserResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "users",
		Verb:     verb,
	}
}
