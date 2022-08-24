/*
Copyright 2021 The Katanomi Authors.

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

package pointer

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBool(t *testing.T) {
	g := NewGomegaWithT(t)
	input := true
	g.Expect(Bool(input)).To(Equal(&input))
}

func TestInt(t *testing.T) {
	g := NewGomegaWithT(t)
	input := int(0)
	g.Expect(Int(input)).To(Equal(&input))
}

func TestInt64(t *testing.T) {
	g := NewGomegaWithT(t)
	input := int64(0)
	g.Expect(Int64(input)).To(Equal(&input))
}

func TestString(t *testing.T) {
	g := NewGomegaWithT(t)
	input := ""
	g.Expect(String(input)).To(Equal(&input))
}

func TestFloat64(t *testing.T) {
	g := NewGomegaWithT(t)
	input := float64(0)
	g.Expect(Float64(input)).To(Equal(&input))
}
