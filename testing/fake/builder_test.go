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

package fake

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/katanomi/pkg/testing/fake/opa"
	. "github.com/onsi/gomega"
)

func TestPolicyBuilderResult(t *testing.T) {
	testCases := map[string]struct {
		result     interface{}
		query      opa.Query
		resultFunc func(policy *opa.Policy) interface{}
	}{
		"map": {
			result: map[string]interface{}{
				"abc": "def",
			},
			resultFunc: func(policy *opa.Policy) interface{} {
				return policy.MapResult("data.fake.result")
			},
		},
		"int": {
			result: 123,
			resultFunc: func(policy *opa.Policy) interface{} {
				return policy.IntResult("data.fake.result")
			},
		},

		"string": {
			result: "123",
			resultFunc: func(policy *opa.Policy) interface{} {
				return policy.StringResult("data.fake.result")
			},
		},
		"object string": {
			result: `{"abc":"def"}`,
			resultFunc: func(policy *opa.Policy) interface{} {
				v := policy.MapResult("data.fake.result")
				result, _ := json.Marshal(v)
				return string(result)
			},
		},
		"bytes": {
			result: []byte("123"),
			resultFunc: func(policy *opa.Policy) interface{} {
				return []byte(policy.StringResult("data.fake.result"))
			},
		},
		"object bytes": {
			result: []byte(`{"abc":"def"}`),
			resultFunc: func(policy *opa.Policy) interface{} {
				v := policy.MapResult("data.fake.result")
				result, _ := json.Marshal(v)
				return result
			},
		},
	}

	for name, item := range testCases {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			builder := NewPolicyBuilder("get", "projects")
			builder.Result(item.result)
			p, err := builder.Complete()
			g.Expect(err).To(BeNil())

			err = p.Eval(context.Background(), nil)
			g.Expect(err).To(BeNil())
			g.Expect(item.resultFunc(p)).To(Equal(item.result))
		})
	}
}

func TestPolicyBuilderWhen(t *testing.T) {
	testCases := map[string]struct {
		expected map[string]interface{}
		input    map[Input]interface{}
		result   interface{}

		field Input
		value interface{}
	}{
		"body": {
			expected: map[string]interface{}{
				"abc": "def",
			},
			input: map[Input]interface{}{
				InputBody: map[string]interface{}{
					"use": 123,
				},
			},
			result: map[string]interface{}{
				"abc": "def",
			},
			field: InputBody.Field("use"),
			value: 123,
		},
		"query": {
			expected: map[string]interface{}{
				"abc": "def",
			},
			input: map[Input]interface{}{
				InputQuery: map[string]interface{}{
					"use": true,
				},
			},
			result: map[string]interface{}{
				"abc": "def",
			},
			field: InputQuery.Field("use"),
			value: true,
		},
		"meta": {
			expected: map[string]interface{}{
				"abc": "def",
			},
			input: map[Input]interface{}{
				InputMeta: map[string]interface{}{
					"baseURL": "http://example.com",
				},
			},
			result: map[string]interface{}{
				"abc": "def",
			},
			field: Input("meta").Field("baseURL"),
			value: "http://example.com",
		},
		"result from input": {
			expected: map[string]interface{}{
				"baseURL": "http://example.com",
			},
			input: map[Input]interface{}{
				"abc": map[string]interface{}{
					"baseURL": "http://example.com",
				},
			},
			result: Input("abc"),
			field:  Input("abc").Field("baseURL"),
			value:  "http://example.com",
		},
		"not matched": {
			expected: map[string]interface{}{},
			input: map[Input]interface{}{
				"asdf": map[string]interface{}{
					"use": "123",
				},
			},
			result: map[string]interface{}{
				"abc": "def",
			},
			field: "qwer.use",
			value: "123",
		},
	}

	for name, item := range testCases {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			builder := NewPolicyBuilder("get", "projects")
			builder.When(item.field, item.value).Result(item.result)
			p, err := builder.Complete()
			g.Expect(err).To(BeNil())

			err = p.Eval(context.Background(), item.input)
			g.Expect(err).To(BeNil())
			result := p.MapResult("data.fake.result")
			g.Expect(result).To(Equal(item.expected))
		})
	}
}
