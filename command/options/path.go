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
	"os/exec"

	"github.com/katanomi/pkg/command/io"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	PathKtnSettingCli = "ktn-settings"
)

// SourcePathOption describe source path option
type SourcePathOption struct {
	SourcePath string
}

// AddFlags add flags to options
func (p *SourcePathOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&p.SourcePath, "source-path", "", `the path contains source code`)
}

// ResultPathOption describe result path option
type ResultPathOption struct {
	// FlagName defines the name when adding the flag
	// defaults to result-path
	FlagName string
	// ResultPath stores the value read from the flag
	ResultPath string
}

// AddFlags add flags to options
func (p *ResultPathOption) AddFlags(flags *pflag.FlagSet) {
	if p.FlagName == "" {
		p.FlagName = "result-path"
	}
	flags.StringVar(&p.ResultPath, p.FlagName, "", `the path to save task results`)
}

// KatanomiPathOption describe katanomi path option
type KatanomiPathOption struct {
	SharePath  string
	BinPath    string
	ConfigPath string
}

// AddFlags add flags to options
func (p *KatanomiPathOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&p.SharePath, "share-path", "/katanomi/data", `the path shared between steps`)
	flags.StringVar(&p.BinPath, "bin-path", "/katanomi/bin", `the path contains binaries`)
	flags.StringVar(&p.ConfigPath, "config-path", "", `the path contains configs`)
}

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
func (p *CLIPathOption) Execute(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, p.CLIPath, args...)
	return cmd.CombinedOutput()
}
