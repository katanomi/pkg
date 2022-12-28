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
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestDependencyReposOption_Setup(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	base := field.NewPath("base")

	obj := struct {
		Cmd DependencyReposOption
	}{Cmd: DependencyReposOption{FlagName: "test-flag"}}
	args := []string{
		"--test-flag", "registry.com", "registry.com",
	}
	err := RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(Succeed(), "parse flag succeed.")
	g.Expect(obj.Cmd.DependencyRepos).To(Equal([]string{"registry.com", "registry.com"}))
	g.Expect(obj.Cmd.Validate(base)).To(HaveLen(0), "validate succeed")
}
