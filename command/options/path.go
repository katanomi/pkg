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
	"github.com/spf13/pflag"
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
	ResultPath string
}

// AddFlags add flags to options
func (p *ResultPathOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&p.ResultPath, "result-path", "", `the path to save task results`)
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
