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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

var _ = Describe("Test.CheckApprovalResults", func() {

	var (
		ctx               context.Context
		approved, denied  bool
		message           string
		newCheckStatus    *metav1alpha1.ApprovalCheckStatus
		approvalPolicy    metav1alpha1.ApprovalPolicy
		approvalSpecUsers []rbacv1.Subject

		approvalSpec     *metav1alpha1.ApprovalSpec
		userApprovals    metav1alpha1.UserApprovals
		approvalStatuses metav1alpha1.ApprovalStatuses

		check       *metav1alpha1.Check
		checkStatus *metav1alpha1.ApprovalCheckStatus

		ApprovedInput = &metav1alpha1.UserApprovalInput{Approved: true}
		DeniedInput   = &metav1alpha1.UserApprovalInput{Approved: false}

		approvalTime = &metav1.Time{Time: time.Now()}
	)

	BeforeEach(func() {
		ctx = logging.WithLogger(context.TODO(), logger)
		approved, denied = false, false
		message = ""
		newCheckStatus = nil

		tm, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
		approvalTime = &metav1.Time{Time: tm}

		approvalPolicy = metav1alpha1.ApprovalPolicyAny
		approvalSpecUsers = []rbacv1.Subject{
			{Name: "user", Kind: rbacv1.UserKind},
			{Name: "admin", Kind: rbacv1.UserKind},
		}
		userApprovals = metav1alpha1.UserApprovals{
			{
				Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
			},
			{
				Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
			},
		}
		approvalStatuses = metav1alpha1.ApprovalStatuses{
			{
				Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
			},
			{
				Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
			},
		}
	})

	JustBeforeEach(func() {
		approvalSpec = &metav1alpha1.ApprovalSpec{
			Policy: approvalPolicy,
			Users:  approvalSpecUsers,
		}
		check = &metav1alpha1.Check{
			Approval: &metav1alpha1.Approval{
				Users: userApprovals,
			},
		}
		checkStatus = &metav1alpha1.ApprovalCheckStatus{
			Approvals: approvalStatuses,
		}
		approved, denied, message, newCheckStatus = CheckApprovalResults(ctx, approvalSpec, check, checkStatus)
	})

	Context("invalid input parameter", func() {
		It("should return an error message", func() {
			By("approval spec is nil")
			_, _, message, _ = CheckApprovalResults(ctx, nil, check, checkStatus)
			Expect(message).Should(Equal(`nil approvalSpec`))
		})
	})

	Context("check status is nil", func() {
		It("should return an not nil check status", func() {
			By("check status is nil")
			_, _, _, newCheckStatus = CheckApprovalResults(ctx, approvalSpec, check, nil)
			Expect(newCheckStatus).ShouldNot(BeNil())
		})
	})

	Context("denied", func() {
		When("only one denied", func() {
			BeforeEach(func() {
				userApprovals[1].Input = DeniedInput
				approvalStatuses[1].UserApprovalInput = DeniedInput
				approvalStatuses[1].ApprovalTime = approvalTime
			})
			It("should denied", func() {
				Expect(denied).To(BeTrue())
				Expect(message).To(Equal("Rejected by \"admin\" on 2006-01-02T15:04:05+07:00"))
			})
		})
	})
	Context("only one approved", func() {
		BeforeEach(func() {
			userApprovals[0].Input = ApprovedInput
			approvalStatuses[0].UserApprovalInput = ApprovedInput
			approvalStatuses[0].ApprovalTime = approvalTime
		})
		When("policy is any", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAny
			})
			It("should passed", func() {
				Expect(approved).To(BeTrue())
				Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00"))
			})
		})
		When("policy is all", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAll
			})
			It("should pending", func() {
				Expect(approved).To(BeFalse())
				Expect(denied).To(BeFalse())
				Expect(message).To(Equal(""))
			})
		})
		When("policy is in order", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyInOrder
			})
			It("should pending", func() {
				Expect(approved).To(BeFalse())
				Expect(denied).To(BeFalse())
				Expect(message).To(Equal(""))
			})
		})
	})

	Context("all approved", func() {
		BeforeEach(func() {
			userApprovals[0].Input = ApprovedInput
			approvalStatuses[0].UserApprovalInput = ApprovedInput
			approvalStatuses[0].ApprovalTime = approvalTime
			userApprovals[1].Input = ApprovedInput
			approvalStatuses[1].UserApprovalInput = ApprovedInput
			approvalStatuses[1].ApprovalTime = approvalTime
		})
		When("policy is any", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAny
			})
			It("should passed", func() {
				Expect(approved).To(BeTrue())
				Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00, \"admin\" on 2006-01-02T15:04:05+07:00"))
			})
		})
		When("policy is all", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAll
			})
			It("should passed", func() {
				Expect(approved).To(BeTrue())
				Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00, \"admin\" on 2006-01-02T15:04:05+07:00"))
			})
		})
		When("policy is in order", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyInOrder
			})
			It("should passed", func() {
				Expect(approved).To(BeTrue())
				Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00, \"admin\" on 2006-01-02T15:04:05+07:00"))
			})
		})
		When("check status is nil", func() {
			It("should return an not nil check status", func() {
				By("check status is nil")
				_, _, _, newCheckStatus = CheckApprovalResults(ctx, approvalSpec, check, nil)
				Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00, \"admin\" on 2006-01-02T15:04:05+07:00"))
				Expect(newCheckStatus).ShouldNot(BeNil())
			})
		})
	})

	Context("append approval user", func() {
		BeforeEach(func() {
			userApprovals = metav1alpha1.UserApprovals{
				{
					Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
				},
				{
					Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
				},
				{
					Subject: rbacv1.Subject{Name: "append", Kind: rbacv1.UserKind},
					Input:   ApprovedInput,
				},
			}
			approvalStatuses = metav1alpha1.ApprovalStatuses{
				{
					Subject: rbacv1.Subject{Name: "user", Kind: rbacv1.UserKind},
				},
				{
					Subject: rbacv1.Subject{Name: "admin", Kind: rbacv1.UserKind},
				},
				{
					Subject:           rbacv1.Subject{Name: "append", Kind: rbacv1.UserKind},
					UserApprovalInput: ApprovedInput,
					ApprovalTime:      approvalTime,
				},
			}
		})
		When("policy is any", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAny
			})
			It("should passed", func() {
				Expect(approved).To(BeTrue())
				Expect(message).To(Equal("Approved by \"append\" on 2006-01-02T15:04:05+07:00"))
			})
		})
		When("policy is all", func() {
			BeforeEach(func() {
				approvalPolicy = metav1alpha1.ApprovalPolicyAll
			})
			When("only two approved", func() {
				BeforeEach(func() {
					userApprovals[0].Input = ApprovedInput
					approvalStatuses[0].UserApprovalInput = ApprovedInput
					approvalStatuses[0].ApprovalTime = approvalTime
				})
				It("should pending", func() {
					Expect(approved).To(BeFalse())
					Expect(denied).To(BeFalse())
					Expect(message).To(Equal(""))
				})
			})
			When("all approved", func() {
				BeforeEach(func() {
					userApprovals[0].Input = ApprovedInput
					approvalStatuses[0].UserApprovalInput = ApprovedInput
					approvalStatuses[0].ApprovalTime = approvalTime
					userApprovals[1].Input = ApprovedInput
					approvalStatuses[1].UserApprovalInput = ApprovedInput
					approvalStatuses[1].ApprovalTime = approvalTime
				})
				It("should passed", func() {
					Expect(approved).To(BeTrue())
					Expect(message).To(Equal("Approved by \"user\" on 2006-01-02T15:04:05+07:00, \"admin\" on 2006-01-02T15:04:05+07:00, \"append\" on 2006-01-02T15:04:05+07:00"))
				})

			})
		})
	})

})
