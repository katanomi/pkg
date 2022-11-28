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
	"strings"

	pkgargs "github.com/katanomi/pkg/command/args"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// KeyValueOptionValidationFunc simple validation function to check the whole
type KeyValueOptionValidationFunc func(path *field.Path, option KeyValueListOption) (errs field.ErrorList)

// KeyValueListOption describe a key=value array set
type KeyValueListOption struct {
	// FlagName used to declare the flag name when interpreting arguments
	// "abc" will expect "--abc key1=value1 key2=value2" argument
	FlagName string
	// KeyValues stored after step
	KeyValues map[string]string

	Validations []KeyValueOptionValidationFunc
}

func (p *KeyValueListOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	p.KeyValues = make(map[string]string)
	flags, _ := pkgargs.GetKeyValues(ctx, args, p.FlagName)

	for k, v := range flags {
		p.KeyValues[strings.Trim(k, ".")] = v
	}
	return nil
}

// AddValidation adds a validation function for this option
func (m *KeyValueListOption) AddValidation(validationFunc ...KeyValueOptionValidationFunc) {
	m.Validations = append(m.Validations, validationFunc...)
}

// Validate validates with all validaton functions
func (m *KeyValueListOption) Validate(path *field.Path) (errs field.ErrorList) {
	for _, val := range m.Validations {
		errs = append(errs, val(path, *m)...)
	}
	return
}

// RequiredKeyValueOptionValidation requires a non-empty values for key and value in all pairs for KeyValueOption
func RequiredKeyValueOptionValidation(path *field.Path, option KeyValueListOption) (errs field.ErrorList) {
	for key, value := range option.KeyValues {
		if strings.TrimSpace(key) == "" {
			errs = append(errs, field.Required(path.Child(key), "key is required."))
		}
		if strings.TrimSpace(value) == "" {
			errs = append(errs, field.Required(path.Child(key), "value is required."))
		}
	}
	return
}

// NotEmptyKeyValueOptionValidation require to have at least one key/value pair
func NotEmptyKeyValueOptionValidation(path *field.Path, option KeyValueListOption) (errs field.ErrorList) {
	if len(option.KeyValues) == 0 {
		errs = append(errs, field.Required(path, "must have key and values."))
	}
	return
}
