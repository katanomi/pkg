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
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	authv1 "k8s.io/api/authorization/v1"
)

var _ = Describe("TestCaseExecution", func() {
	Context("TestCaseExecutionResourceAttributes", func() {
		It("should return related attributes", func() {
			Expect(TestCaseExecutionResourceAttributes("get")).To(Equal(authv1.ResourceAttributes{
				Group:    GroupVersion.Group,
				Version:  GroupVersion.Version,
				Resource: "testcaseexecutions",
				Verb:     "get",
			}))
		})
	})
})

func Test_executorFromNote(t *testing.T) {
	tests := []struct {
		name string
		note string
		want *UserSpec
	}{
		{
			name: "empty string",
			note: "",
			want: nil,
		},
		{
			name: "invalid string format",
			note: "xxxxxx",
			want: nil,
		},
		{
			name: "invalid user string",
			note: "[createdBy: xxx!!!xxx]",
			want: nil,
		},
		{
			name: "valid user string",
			note: "[createdBy: xxx|xxx@xx.x]",
			want: &UserSpec{
				Name:  "xxx",
				Email: "xxx@xx.x",
			},
		},
		{
			name: "allow email is not valid",
			note: "[createdBy: xxx|xxxxx]",
			want: &UserSpec{
				Name:  "xxx",
				Email: "xxxxx",
			},
		},
		{
			name: "allow email is empty",
			note: "[createdBy: xxx|]",
			want: &UserSpec{
				Name:  "xxx",
				Email: "",
			},
		},
		{
			name: "allow both is empty",
			note: "[createdBy: |]",
			want: &UserSpec{
				Name:  "",
				Email: "",
			},
		},
		{
			name: "allow both is empty",
			note: "[createdBy: xiangmu@._-|xiangmu@-._]",
			want: &UserSpec{
				Name:  "xiangmu@._-",
				Email: "xiangmu@-._",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserSpecFromNote(tt.note); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserSpecFromNote() = %v, want %v", got, tt.want)
			}
		})
	}
}
