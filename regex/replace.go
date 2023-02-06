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

import "regexp"

// Replace provide helper functions for replacing strings
type Replace struct {
	// Regex is a regular expression used to modify the original version to generate new variants
	// +optional
	Regex string `json:"regex,omitempty"`

	// Replacement is the value after replacement
	// +optional
	Replacement string `json:"replacement,omitempty"`
}

// Replaces is a list of Replace
type Replaces []Replace

// ReplaceAllString replace all string
func (r *Replace) ReplaceAllString(s string) string {
	re := regexp.MustCompile(r.Regex)
	return re.ReplaceAllString(s, r.Replacement)
}

// ReplaceAllString replace all string
func (rs *Replaces) ReplaceAllString(s string) string {
	if rs == nil {
		return s
	}
	for _, r := range *rs {
		s = r.ReplaceAllString(s)
	}
	return s
}
