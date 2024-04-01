/*
Copyright 2024 The Katanomi Authors.

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

package conformance

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_testCase_AddLabels(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		name   string
		labels []Labels
		want   Labels
	}{
		{
			name:   "empty labels",
			labels: []Labels{},
			want:   nil,
		},
		{
			name:   "multiple labels",
			labels: []Labels{{"test"}, {"abc", "123"}},
			want:   Labels{"test", "abc", "123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			testcase := &testCase{
				node: NewNode(ModuleLevel, "test"),
			}
			testcase.AddLabels(tt.labels...)
			g.Expect(testcase.node.additionalLabels).To(Equal(tt.want))
		})
	}
}
