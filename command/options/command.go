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
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// CommandOption describe command option
type CommandOption struct {
	Command string
}

// AddFlags add flags to options
func (m *CommandOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.Command, "command", "", `command to execute`)
}

// Validate check if command is empty
func (m *CommandOption) Validate(path *field.Path) (errs field.ErrorList) {
	if m.Command == "" {
		errs = append(errs, field.Required(path.Child("command"), "command is required"))
	}
	return errs
}
