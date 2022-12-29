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

	"github.com/katanomi/pkg/command/validators"

	"k8s.io/apimachinery/pkg/util/validation/field"

	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/spf13/cobra"
)

// DependencyReposOption describe dependency repo option
type DependencyReposOption struct {
	DependencyRepos []string
	FlagName        string
}

// Setup init dependency repositories from args
func (m *DependencyReposOption) Setup(ctx context.Context, _ *cobra.Command, args []string) (err error) {
	if m.FlagName == "" {
		m.FlagName = "dependencies-repositories"
	}
	m.DependencyRepos, _ = pkgargs.GetArrayValues(ctx, args, m.FlagName)
	return nil
}

// Validate check if the dependency repository is valid
func (m *DependencyReposOption) Validate(path *field.Path) (errs field.ErrorList) {
	dependencyPath := path.Child(m.FlagName)
	urlValidator := validators.NewURL().SetErrMsg("dependency repository is not a valid url")
	errs = append(errs, urlValidator.Validate(dependencyPath, m.DependencyRepos...)...)
	return errs
}
