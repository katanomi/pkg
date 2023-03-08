/*
Copyright 2023 The Katanomi Authors.

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
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	authv1 "k8s.io/api/authorization/v1"
)

func Test_VariableResourceAttributes(t *testing.T) {
	g := NewGomegaWithT(t)

	want := authv1.ResourceAttributes{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: "variables",
		Verb:     "test",
	}

	got := VariableResourceAttributes("test")
	g.Expect(got).To(Equal(want), "the ResourceAttributes should contain test")
}

func TestVariableList_Filter(t *testing.T) {
	variableList := VariableList{
		Items: []Variable{
			{Name: "var1", Example: "1", Label: "label"},
			{Name: "var2", Example: "1", Label: "label2"},
			{Name: "var3", Example: "2", Label: "label3,label2"},
			{Name: "var4", Example: "3", Label: "label2"},
		},
	}
	tests := map[string]struct {
		filters      []func(*Variable) bool
		variableList VariableList
		want         VariableList
	}{
		"filter not set": {
			variableList: variableList,
			want:         variableList,
		},
		"set label filter": {
			variableList: variableList,
			filters:      []func(*Variable) bool{LabelFilter("label2")},
			want: VariableList{
				Items: []Variable{
					{Name: "var2", Example: "1", Label: "label2"},
					{Name: "var3", Example: "2", Label: "label3,label2"},
					{Name: "var4", Example: "3", Label: "label2"}},
			},
		},
		"set label filter with empty label": {
			variableList: variableList,
			filters:      []func(*Variable) bool{LabelFilter("")},
			want:         VariableList{Items: []Variable{}},
		},
		"set mulit filter": {
			variableList: variableList,
			filters:      []func(*Variable) bool{LabelFilter("label2"), LabelFilter("label3")},
			want: VariableList{
				Items: []Variable{
					{Name: "var3", Example: "2", Label: "label3,label2"},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			tt.variableList.Filter(tt.filters...)
			diff := cmp.Diff(tt.variableList, tt.want)
			g.Expect(diff).To(BeEmpty())
		})
	}
}
