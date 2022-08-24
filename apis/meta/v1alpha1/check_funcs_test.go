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

package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
)

var _ = Describe("Test.Check", func() {
	DescribeTable("Check.GetApprovalUsers",
		func(input *Approval, expected UserApprovals) {
			check := &Check{
				Approval: input,
			}
			actual := check.GetApprovalUsers()
			Expect(actual).To(Equal(expected))
		},
		Entry("nil pointer",
			nil,
			nil,
		),
		Entry("no user",
			&Approval{},
			nil,
		),
		Entry("has user",
			&Approval{
				Users: []UserApproval{{
					Subject: rbacv1.Subject{
						Name: "user1",
					},
				}},
			},
			[]UserApproval{{
				Subject: rbacv1.Subject{
					Name: "user1",
				},
			}},
		),
	)

	DescribeTable("Check.IsNotSet",
		func(input *Approval, expected bool) {
			check := &Check{
				Approval: input,
			}
			actual := check.IsNotSet()
			Expect(actual).To(Equal(expected))
		},
		Entry("nil pointer",
			nil,
			true,
		),
		Entry("no user",
			&Approval{},
			true,
		),
		Entry("has user",
			&Approval{
				Users: []UserApproval{{
					Subject: rbacv1.Subject{
						Name: "user1",
					},
				}},
			},
			false,
		),
	)
})
