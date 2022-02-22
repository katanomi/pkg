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

package matching_test

import (
	. "github.com/katanomi/pkg/user/matching"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
)

var _ = Describe("Match test", func() {

	userConfig := rbacv1.Subject{
		Kind: rbacv1.UserKind,
		Name: "user",
	}

	userWildcard := rbacv1.Subject{
		Kind: rbacv1.UserKind,
		Name: "*",
	}

	groupConfig := rbacv1.Subject{
		Kind: rbacv1.GroupKind,
		Name: "system:masters",
	}

	groupWildcard := rbacv1.Subject{
		Kind: rbacv1.GroupKind,
		Name: "*",
	}

	serviceAccountConfig := rbacv1.Subject{
		Kind:      rbacv1.ServiceAccountKind,
		Name:      "foo",
		Namespace: "bar",
	}

	serviceAccountConfigWildcard := rbacv1.Subject{
		Kind:      rbacv1.ServiceAccountKind,
		Name:      "*",
		Namespace: "bar",
	}

	DescribeTable("#IsRightUser for User",
		func(subject rbacv1.Subject, userName string, matcher gomegatypes.GomegaMatcher) {
			Expect(IsRightUser(authenticationv1.UserInfo{Username: userName}, subject)).To(matcher)
		},
		Entry("no match because request is empty", userConfig, "", BeFalse()),
		Entry("user name is found", userConfig, "user", BeTrue()),
		Entry("user name is not found", userConfig, "user2", BeFalse()),
		Entry("user name is found because of wildcard", userWildcard, "user2", BeTrue()),
	)

	DescribeTable("#IsRightUser for Group",
		func(subject rbacv1.Subject, groupName string, matcher gomegatypes.GomegaMatcher) {
			Expect(IsRightUser(authenticationv1.UserInfo{Groups: []string{groupName}}, subject)).To(matcher)
		},
		Entry("no match because request is empty", groupConfig, "", BeFalse()),
		Entry("group name is found", groupConfig, "system:masters", BeTrue()),
		Entry("group name is not found", groupConfig, "users", BeFalse()),
		Entry("group name is found because of wildcard", groupWildcard, "users", BeTrue()),
	)

	DescribeTable("#IsRightUser for ServiceAccount",
		func(subject rbacv1.Subject, namespace, name string, matcher gomegatypes.GomegaMatcher) {
			Expect(IsRightUser(authenticationv1.UserInfo{
				Username: serviceaccount.MakeUsername(namespace, name),
			}, subject)).To(matcher)
		},
		Entry("no match because request is empty", serviceAccountConfig, "", "", BeFalse()),
		Entry("service account name is found", serviceAccountConfig, "bar", "foo", BeTrue()),
		Entry("service account name is not found", serviceAccountConfig, "bar", "bar", BeFalse()),
		Entry("service account name is found because of wildcard", serviceAccountConfigWildcard, "bar", "users", BeTrue()),
		Entry("service account name is found because of different namespace", serviceAccountConfigWildcard, "foo", "foo", BeFalse()),
	)
})

var _ = Describe("Convert test", func() {
	var (
		defaultNamespace = "default"
		userInfo         authenticationv1.UserInfo
		subject          rbacv1.Subject
	)

	JustBeforeEach(func() {
		subject = ConvertUserInfoToSubject(userInfo, defaultNamespace)
	})

	When("kind is user", func() {
		BeforeEach(func() {
			userInfo.Username = "user"
		})
		It("should convert to user", func() {
			expectUser := rbacv1.Subject{
				Kind: rbacv1.UserKind,
				Name: "user",
			}
			Expect(subject).To(Equal(expectUser))
		})
	})

	When("kind is sa", func() {
		BeforeEach(func() {
			userInfo.Username = "system:serviceaccount:default"
		})
		It("should convert to sa", func() {
			expectSa := rbacv1.Subject{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      "default",
				Namespace: defaultNamespace,
			}
			Expect(subject).To(Equal(expectSa))
		})
	})
})
