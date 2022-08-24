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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test.Approval", func() {
	DescribeTable("Approval.HasApprover",
		func(input *Approval, expected bool) {
			actual := input.HasApprover()
			Expect(actual).To(Equal(expected))
		},
		Entry("nil pointer",
			nil,
			false,
		),
		Entry("no user",
			&Approval{},
			false,
		),
		Entry("has user",
			&Approval{
				Users: []UserApproval{{
					Subject: rbacv1.Subject{
						Name: "user1",
					},
				}},
			},
			true,
		),
	)
})

var _ = Describe("Test.ApprovalSpec", func() {
	DescribeTable("ApprovalSpec.HasApprover.IsEnabled.TimeoutEnabled",
		func(input *ApprovalSpec, expected bool) {
			var actual bool
			actual = input.HasApprover()
			Expect(actual).To(Equal(expected))
			actual = input.IsEnabled()
			Expect(actual).To(Equal(expected))
			actual = input.TimeoutEnabled()
			Expect(actual).To(Equal(expected))
		},
		Entry("nil pointer",
			nil,
			false,
		),
		Entry("no user",
			&ApprovalSpec{},
			false,
		),
		Entry("has user",
			&ApprovalSpec{
				Users: []rbacv1.Subject{{
					Name: "user1",
				}},
				Timeout: metav1.Duration{
					Duration: 1,
				},
			},
			true,
		),
	)
})
