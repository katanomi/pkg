/*
Copyright 2021 The Katanomi Authors.

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

package validation

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateName standard name validation for all katanomi types
// uses the same validation of namespace objects
func ValidateName(obj metav1.Object, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if fldPath == nil {
		fldPath = field.NewPath("metadata", "name")
	}
	errList := apimachineryvalidation.ValidateNamespaceName(obj.GetName(), false)
	if len(errList) > 0 {
		for _, errStr := range errList {
			allErrs = append(allErrs, field.Invalid(fldPath, obj.GetName(), errStr))
		}
	}
	return allErrs
}

// ValidateAnnotations to contain size and data format
var ValidateAnnotations = apimachineryvalidation.ValidateAnnotations
