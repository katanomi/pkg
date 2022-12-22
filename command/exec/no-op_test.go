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

func TestNoOpCmder(t *testing.T) {
	ctx := context.TODO()
	cmder := &NoOpCmder{}
	t.Run("no op cmdr Command", func(t *testing.T) {
		g := NewGomegaWithT(t)

		cmd := cmder.Command("command", "arg1", "arg2")
		cmd.SetEnv("key=value")
		cmd.SetStdin(nil)
		lines, err := CombinedOutputLines(cmd)
		g.Expect(err).To(BeNil())
		g.Expect(lines).To(Equal([]string{"command arg1 arg2"}))
	})

	t.Run("no op cmdr CommandContext", func(t *testing.T) {
		g := NewGomegaWithT(t)

		cmd := cmder.CommandContext(ctx, "command1", "argX", "argZ")
		lines, err := CombinedOutputLines(cmd)
		g.Expect(err).To(BeNil())
		g.Expect(lines).To(Equal([]string{"command1 argX argZ"}))
	})
}
