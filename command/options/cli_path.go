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

	"github.com/katanomi/pkg/command/exec"
	"github.com/katanomi/pkg/command/io"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// CLIPathOption adds a generic option to store different cli paths
type CLIPathOption struct {
	// Name of the program
	// i.e yq, helm etc.
	Name string
	// CLIPath direct path for CLI
	// used to store the default value
	CLIPath string
	// FlagName to be used to store
	FlagName string
}

// AddFlags adds flags for option
func (p *CLIPathOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&p.CLIPath, p.FlagName, p.CLIPath, `the path for `+p.Name)
}

// Validate if values are given
func (p *CLIPathOption) Validate(path *field.Path) (errs field.ErrorList) {
	if p.CLIPath == "" {
		errs = append(errs, field.Required(path.Child(p.FlagName), `path for `+p.Name+` is necessary`))
	} else if !io.IsExist(p.CLIPath) {
		errs = append(errs, field.Required(path.Child(p.FlagName), `path "`+p.CLIPath+`" for `+p.Name+` does not exist`))
	}
	return
}

// Execute executes code given a context
func (p *CLIPathOption) Execute(ctx context.Context, args ...string) ([]string, error) {
	cmder := exec.FromContextCmder(ctx)
	cmd := cmder.CommandContext(ctx, p.CLIPath, args...)
	return exec.CombinedOutputLines(cmd)
}
