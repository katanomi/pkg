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

package assertions

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var stringPtr *string

var nonEmptyString = "abc"
var emptyString = ""
var nonEmptyInt = 1
var emptyInt = 0

type SomeType struct {
	A string
}

var emptyType = SomeType{}
var nonEmptyType = SomeType{A: "a"}
var typePointer *SomeType

var _ = DescribeTable("BeNilOrEmpty",
	func(value any, expected bool, expectedErr error) {
		success, err := BeNilOrEmpty().Match(value)
		Expect(success).To(Equal(expected))
		if expectedErr == nil {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
	Entry("nil", nil, true, nil),
	// string
	Entry("empty string", "", true, nil),
	Entry("non empty string", nonEmptyString, false, nil),
	Entry("nil string pointer", stringPtr, true, nil),
	Entry("non nil empty string pointer", &emptyString, true, nil),
	Entry("non empty string pointer", &nonEmptyString, false, nil),
	// int
	Entry("zero int", 0, true, nil),
	Entry("non zero int", 1, false, nil),
	Entry("non nil empty int pointer", &emptyInt, true, nil),
	Entry("non empty int pointer", &nonEmptyInt, false, nil),
	// structs
	Entry("nil pointer", typePointer, true, nil),
	Entry("non nil empty struct", emptyType, true, nil),
	Entry("non nil empty struct pointer", &emptyType, true, nil),
	Entry("non empty struct pointer", &nonEmptyType, false, nil),
)
