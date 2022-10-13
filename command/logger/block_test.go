/*
Copyright 2022 The Katanomi Authors.

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

package logger

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBlock(t *testing.T) {
	g := NewGomegaWithT(t)
	out := []string{}
	ctx := genTestLogContext(&out)

	Block(ctx, "test", "test content")
	g.Expect(out).To(HaveLen(1))
	g.Expect(out[0]).To(Equal(`
====================== test ======================
test content
==================================================
`,
	))

	Block(ctx, "", "")
	g.Expect(out).To(HaveLen(2))
	g.Expect(out[1]).To(Equal(`
========================  ========================

==================================================
`,
	))

	longStr := "123456789012345678901234567890123456789012345678901234567890"
	Block(ctx, longStr, longStr)
	g.Expect(out).To(HaveLen(3))
	g.Expect(out[2]).To(Equal(`
= 123456789012345678901234567890123456789012345678901234567890 =
123456789012345678901234567890123456789012345678901234567890
==================================================
`,
	))
}
