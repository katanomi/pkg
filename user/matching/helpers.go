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

package matching

import (
	"fmt"
	"strings"

	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
)

// WildcardAll is a character which represents all elements in a set.
const WildcardAll = "*"

// UserMatches returns `true` if the given user in the subject has a match in the given userConfig.
func UserMatches(subject rbacv1.Subject, userInfo authenticationv1.UserInfo) bool {
	if subject.Kind != rbacv1.UserKind {
		return false
	}

	return subject.Name == WildcardAll || subject.Name == userInfo.Username
}

// UserGroupMatches returns `true` if the given group in the subject has a match in the given userConfig.
// Always returns true if `WildcardAll` is used in subject.
func UserGroupMatches(subject rbacv1.Subject, userInfo authenticationv1.UserInfo) bool {
	if subject.Kind != rbacv1.GroupKind {
		return false
	}

	if subject.Name == WildcardAll {
		return true
	}

	for _, group := range userInfo.Groups {
		if group == subject.Name {
			return true
		}
	}
	return false
}

// ServiceAccountMatches returns `true` if the given service account in the subject has a match in the given userConfig.
// Supports `WildcardAll` in subject name.
func ServiceAccountMatches(subject rbacv1.Subject, userInfo authenticationv1.UserInfo) bool {
	if subject.Kind != rbacv1.ServiceAccountKind {
		return false
	}

	if subject.Name == WildcardAll {
		saPrefix := fmt.Sprintf("%s%s:", serviceaccount.ServiceAccountUsernamePrefix, subject.Namespace)
		return strings.HasPrefix(userInfo.Username, saPrefix)
	}

	return serviceaccount.MatchesUsername(subject.Namespace, subject.Name, userInfo.Username)
}
