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
	"github.com/katanomi/pkg/command/io"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ContextOption describe the work dir to execute the command
type ContextOption struct {
	Context               string
	ValidateContextExists bool
}

// AddFlags add flags to options
func (p *ContextOption) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&p.Context, "context", "", `the work dir to execute the command`)
}

// Validate check if command is empty
func (p *ContextOption) Validate(path *field.Path) (errs field.ErrorList) {
	if p.Context == "" {
		errs = append(errs, field.Required(path.Child("context"), "context is required"))
	}
	if p.ValidateContextExists && !io.IsDir(p.Context) {
		errs = append(errs, field.Invalid(path.Child("context"), p.Context, `context is not a folder`))
	}
	return errs
}
