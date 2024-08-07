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

package script

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestExecScript(t *testing.T) {
	g := NewGomegaWithT(t)

	scriptResult := ExecScript("echo", "happy test")

	g.Expect(scriptResult.Stderr()).To(BeEmpty())
	g.Expect(scriptResult.Stdout()).NotTo(BeEmpty())
	g.Expect(scriptResult.Error()).To(BeNil())
	g.Expect(scriptResult.ExitMessage()).To(BeEmpty())
	g.Expect(scriptResult.ExitCode()).To(Equal(0))

	scriptResult = ExecScript("cat", "not_found_file.tt")

	g.Expect(scriptResult.Stderr()).NotTo(BeEmpty())
	g.Expect(scriptResult.Stdout()).To(BeEmpty())
	g.Expect(scriptResult.Error()).NotTo(BeNil())
	g.Expect(scriptResult.ExitMessage()).NotTo(BeEmpty())
	g.Expect(scriptResult.ExitCode()).NotTo(Equal(0))
}
