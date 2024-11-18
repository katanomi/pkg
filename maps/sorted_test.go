/*
Copyright 2022 The AlaudaDevops Authors.

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

package maps

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSortedKeyValue(t *testing.T) {

	table := map[string]struct {
		Entry  map[string]string
		Result []KeyValue
	}{
		"multiple keys and values": {
			Entry: map[string]string{
				"a":  "bc",
				"b":  "11dex1",
				"ab": "1231",
			},
			Result: []KeyValue{
				{"a", "bc"},
				{"ab", "1231"},
				{"b", "11dex1"},
			},
		},
		"empty": {
			Entry:  map[string]string{},
			Result: []KeyValue{},
		},
		"nil": {
			Entry:  nil,
			Result: []KeyValue{},
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := SortedKeyValue(test.Entry)
			g.Expect(result).To(Equal(test.Result))
		})
	}
}
