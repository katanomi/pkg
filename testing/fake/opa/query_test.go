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

package opa

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestQueryVar(t *testing.T) {
	q := Query("a.b.c")
	g := NewGomegaWithT(t)
	g.Expect(q.Var()).To(Equal("a_b_c"))
}

func TestQueryEval(t *testing.T) {
	q := Query("a.b.c")
	g := NewGomegaWithT(t)
	g.Expect(q.Eval()).To(Equal("a_b_c = a.b.c"))
}

func TestQueriesEval(t *testing.T) {
	q := Queries{"a.b.c", "x.y.z"}
	g := NewGomegaWithT(t)
	g.Expect(q.Eval()).To(Equal("a_b_c = a.b.c;x_y_z = x.y.z"))
}
