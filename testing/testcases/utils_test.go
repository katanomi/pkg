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

package testcases

import (
	"testing"

	. "github.com/onsi/gomega"
)

type mockNameGetter struct {
	name string
}

func (m *mockNameGetter) GetName() string {
	return m.name
}

func TestFindByName(t *testing.T) {
	g := NewGomegaWithT(t)
	var (
		t1 = &mockNameGetter{name: "foo"}
		t2 = &mockNameGetter{name: "bar"}
		t3 = &mockNameGetter{name: "baz"}
	)

	testCases := map[string]struct {
		list []NameGetter
		name string
		want NameGetter
	}{
		"should return nil for empty list": {
			list: []NameGetter{},
			name: "foo",
			want: nil,
		},
		"should find item by name": {
			list: []NameGetter{t1, t2, t3},
			name: "bar",
			want: t2,
		},
		"should return nil for non-existent name": {
			list: []NameGetter{t1, t2, t3},
			name: "qux",
			want: nil,
		},
	}

	for testName, testCase := range testCases {
		got := FindByName(testCase.list, testCase.name)
		if testCase.want == nil {
			g.Expect(got).To(BeNil(), testName)
		} else {
			g.Expect(got).ToNot(BeNil(), testName)
			g.Expect(got.GetName()).To(Equal(testCase.want.GetName()), testName)
		}
	}
}
