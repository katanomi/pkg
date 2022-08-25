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

package approval

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"knative.dev/pkg/logging"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

var _ = Describe("Test.ValidateApproval", func() {

	var (
		reqUser              authenticationv1.UserInfo
		allowRepresentOthers bool
		checkList            []PairOfOldNewCheck
		approvalSpecList     []*metav1alpha1.ApprovalSpec

		err               error
		username          string
		old, new          metav1alpha1.UserApprovals
		approvalPolicy    metav1alpha1.ApprovalPolicy
		approvalSpecUsers []rbacv1.Subject
	)

	BeforeEach(func() {
		approvalPolicy = metav1alpha1.ApprovalPolicyAny
		username = "admin"
		allowRepresentOthers = false
		approvalSpecUsers = []rbacv1.Subject{}
	})

	JustBeforeEach(func() {
		ctx := logging.WithLogger(context.TODO(), logger)
		reqUser = authenticationv1.UserInfo{
			Username: username,
		}
		approvalSpecList = []*metav1alpha1.ApprovalSpec{
			{
				Policy: approvalPolicy,
				Users:  approvalSpecUsers,
			},
		}
		checkList = []PairOfOldNewCheck{{
			&metav1alpha1.Check{Approval: &metav1alpha1.Approval{Users: old}},
			&metav1alpha1.Check{Approval: &metav1alpha1.Approval{Users: new}},
		}}
		err = ValidateApproval(ctx, reqUser, allowRepresentOthers, approvalSpecList, checkList)
	})

	Context("invalid input parameter", func() {
		It("should return an error", func() {
			By("approval spec list is nil")
			err = ValidateApproval(context.TODO(), reqUser, allowRepresentOthers, nil, checkList)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(BeEquivalentTo(`internal error #check != #checkSpec`))
		})
	})

	Context("allow to represent others", func() {
		BeforeEach(func() {
			username = "admin"
			allowRepresentOthers = true
		})

		When("there is no change in approval", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("append approval", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("approves directly", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("approves for others", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
					{
						Subject: rbacv1.Subject{Name: "user-1", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("approver repeated", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: true},
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: true},
					},
				}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`approver "user" cannot be repeated`))
			})
		})

		When("change the approval result", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: true},
					},
				}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`unable to change the approval result for "user" from &{false } to &{true }`))
			})
		})

		When("remove approver from list", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
				new = []metav1alpha1.UserApproval{{}}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`cannot remove "user" from the approval list`))
			})
		})

	})

	Context("cannot to represent others", func() {
		BeforeEach(func() {
			username = "user"
			allowRepresentOthers = false
		})

		When("approves directly", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("append user but not in spec", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`"user" can not change the approval user list`))
			})
		})

		When("append user but in spec", func() {
			BeforeEach(func() {
				approvalSpecUsers = []rbacv1.Subject{{Name: "user", Kind: rbacv1.UserKind}}
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("approve for other", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: false},
					},
				}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`"user" can not approve for user "admin"`))
			})
		})

	})

	Context("approval policy is InOrder", func() {
		BeforeEach(func() {
			username = "admin"
			allowRepresentOthers = true
			approvalPolicy = metav1alpha1.ApprovalPolicyInOrder
		})

		When("approves in order", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: true},
					},
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
				}
			})
			It("should pass", func() {
				Expect(err).To(BeNil())
			})
		})

		When("approves out of order", func() {
			BeforeEach(func() {
				old = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   nil,
					},
				}
				new = []metav1alpha1.UserApproval{
					{
						Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
					},
					{
						Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
						Input:   &metav1alpha1.UserApprovalInput{Approved: true},
					},
				}
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(BeEquivalentTo(`Approval policy is "InOrder", "user" can not approve before "admin".`))
			})
		})

	})

})
