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

	"github.com/golang-jwt/jwt/v4"
	"github.com/katanomi/pkg/substitution"
	authenticationv1 "k8s.io/api/authentication/v1"
	authv1 "k8s.io/api/authorization/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var UserGVK = GroupVersion.WithKind("User")
var UserListGVK = GroupVersion.WithKind("UserList")

// User object for plugin
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec"`
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
		Resource: "virtualusers",
		Verb:     verb,
	}
}

// UserInfo is generic user information types
type UserInfo struct {
	// UserInfo holds the information about the user needed to implement the user.Info interface.
	authenticationv1.UserInfo
}

const (
	userNameKey  string = "name"
	userEmailKey string = "email"
)

// FromJWT get information from jwt to populate userinfo
func (user *UserInfo) FromJWT(claims jwt.MapClaims) {

	_username, ok := claims[userNameKey]
	if ok {
		username := _username.(string)
		user.Username = username
	}
	_email, ok := claims[userEmailKey]
	if ok {
		email := _email.(string)
		user.Extra = make(map[string]authenticationv1.ExtraValue)
		user.Extra[userEmailKey] = authenticationv1.ExtraValue{email}
	}
	return
}

// GetName get username from UserInfo
func (user *UserInfo) GetName() string {
	return user.Username
}

// GetEmail get email from UserInfo
func (user *UserInfo) GetEmail() string {
	extraValue, ok := user.Extra[userEmailKey]
	if !ok {
		return ""
	}
	if len(extraValue) > 0 {
		return extraValue[0]
	}
	return ""
}

// RBACSubjectValGetter returns the list of keys and values to support variable substitution for
// rbac.
func RBACSubjectValGetter(subject *rbac.Subject) substitution.GetValWithKeyFunc {
	if subject == nil {
		subject = &rbac.Subject{}
	}
	return func(ctx context.Context, path *field.Path) (values map[string]string) {
		values = map[string]string{
			path.String():                    "",
			path.Child("kind").String():      subject.Kind,
			path.Child("apiGroup").String():  subject.APIGroup,
			path.Child("name").String():      subject.Name,
			path.Child("namespace").String(): subject.Namespace,
		}
		return
	}
}
