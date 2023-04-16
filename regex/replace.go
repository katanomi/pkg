/*
Copyright 2023 The Katanomi Authors.

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

package regex

import (
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Replace provides helper functions for replacing strings
type Replace struct {
	// Regex is a regular expression used to modify the original version to generate new variants
	// +optional
	Regex string `json:"regex,omitempty"`

	// Replacement is the value after replacement
	// +optional
	Replacement string `json:"replacement,omitempty"`

	// ToLower is a flag to convert the string to lowercase
	ToLower bool `json:"toLower,omitempty"`

	// ToUpper is a flag to convert the string to uppercase
	ToUpper bool `json:"toUpper,omitempty"`
}

// Replaces is a list of Replace
type Replaces []Replace

// ReplaceAllString replace all string
func (r *Replace) ReplaceAllString(s string) string {
	if r == nil || r.Regex == "" {
		return s
	}
	re := regexp.MustCompile(r.Regex)
	s = re.ReplaceAllString(s, r.Replacement)
	if r.ToLower {
		s = strings.ToLower(s)
	}
	if r.ToUpper {
		s = strings.ToUpper(s)
	}
	return s
}

// ReplaceAllString replace all string
func (rs *Replaces) ReplaceAllString(s string) string {
	if rs == nil || len(*rs) == 0 {
		return s
	}
	for _, r := range *rs {
		s = r.ReplaceAllString(s)
	}
	return s
}

// Validate Replace validation method
func (r *Replace) Validate(fld *field.Path) (errs field.ErrorList) {
	if r.Regex == "" {
		return
	}
	_, err := regexp.Compile(r.Regex)
	if err != nil {
		errs = append(errs, field.Invalid(fld.Child("regex"), r.Regex, err.Error()))
	}

	if r.ToLower && r.ToUpper {
		errs = append(errs, field.Invalid(fld, r.ToLower, "toLower and toUpper cannot be set at the same time"))
	}

	return
}

// Validate Replaces validation method
func (r *Replaces) Validate(fld *field.Path) (errs field.ErrorList) {
	if r == nil {
		return
	}
	for i, replace := range *r {
		errs = append(errs, replace.Validate(fld.Index(i))...)
	}

	return
}
