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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var _ = Describe("Test.IsTimeout", func() {
	DescribeTable("IsTimeout",
		func(startTime *metav1.Time, timeout time.Duration, expected bool) {
			v1time := metav1.Duration{Duration: timeout}
			actual := IsTimeout(startTime, v1time)
			Expect(actual).To(Equal(expected))
		},
		Entry("nil startTime",
			nil,
			1*time.Second,
			false,
		),
		Entry("timeout is zero",
			&metav1.Time{Time: time.Now()},
			0*time.Second,
			false,
		),
		Entry("not timeout",
			&metav1.Time{Time: time.Now().Add(1 * time.Minute)},
			30*time.Second,
			false,
		),
		Entry("timeout",
			&metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
			30*time.Second,
			true,
		),
		Entry("just a timeout",
			&metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
			60*time.Second,
			true,
		),
	)
})
