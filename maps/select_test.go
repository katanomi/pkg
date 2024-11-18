/*
Copyright 2024 The AlaudaDevops Authors.

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

// Package maps contains methods to operate maps of all kinds
package maps

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var source = map[string]string{
	"core-1": "1",
	"core-2": "2",
	"core-3": "3",
}
var _ = Describe("SelectAndMutateMap", func() {
	DescribeTable("SelectAndMutateMap in different situation",
		func(source map[string]string, match MatchFunc, mutate MutateFunc, expect map[string]string) {
			result := SelectAndMutateMap(source, match, mutate)
			Expect(result).To(Equal(expect))
		},
		Entry("return all the element",
			source,
			func(key string, value string) bool { return true },
			func(key string, value string) (string, string) { return key, value },
			source,
		),
		Entry("return core-1 element",
			source,
			func(key string, value string) bool {
				return key == "core-1"
			},
			func(key string, value string) (string, string) { return key, value },
			map[string]string{
				"core-1": "1",
			},
		),
		Entry("return mutated core-1 element",
			source,
			func(key string, value string) bool {
				return key == "core-1"
			},
			func(key string, value string) (string, string) {
				return strings.TrimPrefix(key, "core-"), value
			},
			map[string]string{
				"1": "1",
			},
		),
	)
})
