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

package options

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestEnvFlagsOption(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	obj := struct {
		EnvFlagsOption
	}{}
	args := []string{
		"--env-flags", "FOO=BAR",
	}
	err := RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(Succeed())
	g.Expect(obj.EnvFlags).To(Equal(map[string]string{
		"FOO": "BAR",
	}))

	// empty args
	emptyArgs := []string{}
	g.Expect(RegisterSetup(&obj, ctx, nil, emptyArgs)).ShouldNot(HaveOccurred())

	// invalid args
	invalidArgs := []string{
		"--env-flags", "XXXX",
	}
	g.Expect(RegisterSetup(&obj, ctx, nil, invalidArgs)).Should(HaveOccurred())

}
