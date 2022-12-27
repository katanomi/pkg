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

package exec

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCmderContext(t *testing.T) {
	ctx := context.TODO()
	cmder := &LocalCmder{}
	t.Run("empty cmder from context returns DefaultCmder", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(FromContextCmder(ctx)).To(Equal(DefaultCmder))
	})

	t.Run("add cmder into context and extracts", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctx = WithCmder(ctx, cmder)
		g.Expect(FromContextCmder(ctx)).To(Equal(cmder))
	})
}

func TestCmdContext(t *testing.T) {
	ctx := context.TODO()
	cmd := &LocalCmd{}
	t.Run("empty cmd from context", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(FromContextCmd(ctx)).To(BeNil())
	})

	t.Run("add cmder into context and extracts", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctx = WithCmd(ctx, cmd)
		g.Expect(FromContextCmd(ctx)).To(Equal(cmd))
	})
}
