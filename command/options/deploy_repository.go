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

	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/katanomi/pkg/command/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DeployRepositoryOption describe deploy repo option
type DeployRepositoryOption struct {
	// DeployRepository deploy registry url
	DeployRepository string

	// Required define DeployRepository is required
	Required bool

	// DeployArgs save deloy agrs.
	DeployArgs map[string]string
}

// Setup init deploy args from args
func (m *DeployRepositoryOption) Setup(ctx context.Context, _ *cobra.Command, args []string) (err error) {
	m.DeployArgs, err = pkgargs.GetKeyValues(ctx, args, "deploy-args")
	return err
}

// AddFlags add flags to options
func (m *DeployRepositoryOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.DeployRepository, "deploy-repository", "", `deploy artifact repository url.`)
}

// Validate check if the deploy repository is valid
func (m *DeployRepositoryOption) Validate(path *field.Path) (errs field.ErrorList) {
	deployPath := path.Child("deploy-repository")
	if m.Required && m.DeployRepository == "" {
		errs = append(errs, field.Required(deployPath, "deploy repository must be set"))
		return
	}

	urlValidator := validators.NewURL().SetErrMsg("deploy repository is not a valid url")
	errs = append(errs, urlValidator.Validate(deployPath, m.DeployRepository)...)
	return errs
}
