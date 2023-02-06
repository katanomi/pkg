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
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

func TestUnmarshalAggregateResult(t *testing.T) {
	ret := AggregateResult{
		{
			"revision": "123",
		},
	}
	type Result struct {
		Revision string `json:"revision"`
	}
	retList := []Result{}
	err := ret.Unmarshal(&retList)
	fmt.Println(err, retList)
}

func TestAggregateResult_Unmarshal(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	aggRet := AggregateResult{
		{
			"name": "123",
			"age":  12,
		},
		{
			"name": "456",
			"age":  13,
		},
	}
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	tests := map[string]struct {
		list       interface{}
		wantErr    bool
		wantResult interface{}
	}{
		"unmarshal to struct": {
			list:    new([]Person),
			wantErr: false,
			wantResult: &[]Person{
				{Name: "123", Age: 12},
				{Name: "456", Age: 13},
			},
		},
		"type mismatch": {
			list:    &[]int{},
			wantErr: true,
		},
		"list param is not a slice pointer": {
			list:    Person{},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := aggRet.Unmarshal(tt.list)
			if tt.wantErr {
				g.Expect(err).NotTo(gomega.BeNil())
			} else {
				g.Expect(tt.list).To(gomega.Equal(tt.wantResult))
			}
		})
	}
}

// . "github.com/onsi/gomega"
func TestAggregateFunc(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	g.Expect(Max("test-key", "test-key-alias")).To(gomega.Equal(AggregateField{
		Field:    Field{Name: "test-key", Alias: "test-key-alias"},
		Operator: AggregateOperatorMax,
	}))
	g.Expect(Min("test-key", "test-key-alias")).To(gomega.Equal(AggregateField{
		Field:    Field{Name: "test-key", Alias: "test-key-alias"},
		Operator: AggregateOperatorMin,
	}))
	g.Expect(Count("test-key-alias")).To(gomega.Equal(AggregateField{
		Field:    Field{Name: "", Alias: "test-key-alias"},
		Operator: AggregateOperatorCount,
	}))
	g.Expect(Sum("test-key", "test-key-alias")).To(gomega.Equal(AggregateField{
		Field:    Field{Name: "test-key", Alias: "test-key-alias"},
		Operator: AggregateOperatorSum,
	}))
}
