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

// ResultsOCIContainerImageBuildURLOption describe ociContainerImageBuild-url option
// Deprecated: use ContainerImageOption result store mechanism instead
type ResultsOCIContainerImageBuildURLOption struct {
	OciContainerImageBuildUrl string
}

// AddFlags add flags to options
func (m *ResultsOCIContainerImageBuildURLOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.OciContainerImageBuildUrl, "results-ociContainerImageBuild-url", "", `the path to save ociContainerImageBuild-url, details: https://github.com/katanomi/spec/blob/main/docs/core/contracts/3.core.task-contracts.ociContainerImageBuild.md#results`)
}
