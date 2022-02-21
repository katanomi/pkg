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
	"strings"

	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
)

func userMatch(userInfo authenticationv1.UserInfo, subjects []rbacv1.Subject) bool {
	for _, subject := range subjects {
		var match bool
		switch subject.Kind {
		case rbacv1.UserKind:
			match = UserMatches(subject, userInfo)
		case rbacv1.GroupKind:
			match = UserGroupMatches(subject, userInfo)
		}
		if match {
			return true
		}
	}
	return false
}

func serviceAccountMatch(userInfo authenticationv1.UserInfo, subjects []rbacv1.Subject) bool {
	for _, subject := range subjects {
		if subject.Kind == rbacv1.ServiceAccountKind {
			if ServiceAccountMatches(subject, userInfo) {
				return true
			}
		}
	}
	return false
}

// IsRightUser determine whether the two types of users match
func IsRightUser(userInfo authenticationv1.UserInfo, subject rbacv1.Subject) bool {
	subjects := []rbacv1.Subject{subject}
	isServiceAccount := strings.HasPrefix(userInfo.Username, serviceaccount.ServiceAccountUsernamePrefix)
	if isServiceAccount {
		return serviceAccountMatch(userInfo, subjects)
	}
	return userMatch(userInfo, subjects)
}

func ConvertUserInfoToSubject(userInfo authenticationv1.UserInfo, namespace string) (subject rbacv1.Subject) {
	isServiceAccount := strings.HasPrefix(userInfo.Username, serviceaccount.ServiceAccountUsernamePrefix)
	if isServiceAccount {
		return rbacv1.Subject{
			Kind:      rbacv1.ServiceAccountKind,
			Name:      strings.TrimPrefix(userInfo.Username, serviceaccount.ServiceAccountUsernamePrefix),
			Namespace: namespace,
		}
	}

	return rbacv1.Subject{
		Kind: rbacv1.UserKind,
		Name: userInfo.Username,
	}
}
