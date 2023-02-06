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

package variable

import (
	"strings"

	"k8s.io/utils/strings/slices"
)

// LabelFilter label filter func
func LabelFilter(label string) FilterFunc {
	return func(variable *Variable) bool {
		if variable == nil {
			return false
		}

		if label == "" {
			return true
		}

		varLabels := strings.Split(variable.Label, ",")
		return slices.Contains(varLabels, label)
	}
}

func filtVariable(s *Variable, filters ...FilterFunc) bool {
	for _, f := range filters {
		if f != nil && f(s) == false {
			return false
		}
	}
	return true
}
