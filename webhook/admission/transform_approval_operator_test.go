/*
Copyright 2022 The Katanomi Authors.

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

package admission

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/namespace"
)

var _ = Describe("Test.addApprovalOperator", func() {

	var (
		username         string
		defaultNamespace = "default"
		reqUser          authenticationv1.UserInfo
		checkList        []PairOfOldNewCheck
		old, new         metav1alpha1.UserApprovals
		ctx              context.Context
	)

	JustBeforeEach(func() {
		ctx = namespace.WithNamespace(context.Background(), defaultNamespace)
		reqUser = authenticationv1.UserInfo{
			Username: username,
		}
		checkList = []PairOfOldNewCheck{{
			&metav1alpha1.Check{Approval: &metav1alpha1.Approval{Users: old}},
			&metav1alpha1.Check{Approval: &metav1alpha1.Approval{Users: new}},
		}}
		addApprovalOperator(ctx, reqUser, checkList)
	})

	Context("Approval for myself", func() {
		When("old user is empty", func() {
			BeforeEach(func() {
				username = "user"
				old = nil
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
						Operator: &rbacv1.Subject{
							Name: "user",
							Kind: rbacv1.UserKind,
						},
					},
				}
			})
			It("should be clear the operator", func() {
				Expect(new[0].Operator).To(BeNil())
			})
		})
		When("exists the operator", func() {
			BeforeEach(func() {
				username = "user"
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
						Operator: &rbacv1.Subject{
							Name: "user",
							Kind: rbacv1.UserKind,
						},
					},
				}
			})
			It("should be clear the operator", func() {
				Expect(new[0].Operator).To(BeNil())
			})
		})
	})

	Context("Approval for someone else", func() {
		When("operator is user", func() {
			BeforeEach(func() {
				username = "admin"
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should add the operator", func() {
				expectUser := &rbacv1.Subject{
					Kind: rbacv1.UserKind,
					Name: "admin",
				}
				Expect(new[0].Operator).To(Equal(expectUser))
			})
		})

		When("operator is sa", func() {
			BeforeEach(func() {
				username = "system:serviceaccount:default"
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should add the operator", func() {
				expectSa := &rbacv1.Subject{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      "default",
					Namespace: defaultNamespace,
				}
				Expect(new[0].Operator).To(Equal(expectSa))
			})
		})
	})
})
