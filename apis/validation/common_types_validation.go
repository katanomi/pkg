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
	"k8s.io/apimachinery/pkg/runtime/schema"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateName standard name validation for all katanomi types
// uses the same validation of namespace objects
func ValidateName(obj metav1.Object, fld *field.Path) field.ErrorList {
	return ValidateItemName(obj.GetName(), false, fld)
}

// ValidateAnnotations to contain size and data format
var ValidateAnnotations = apimachineryvalidation.ValidateAnnotations

// ValidateObjectReference validates an object reference
func ValidateObjectReference(objref *corev1.ObjectReference, optional, needsResourceType bool, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	if objref == nil && !optional {
		errs = append(errs, field.Required(fld, "a valid reference is required"))
	}
	if objref != nil {
		if objref.Kind == "" && needsResourceType {
			errs = append(errs, field.Required(fld.Child("kind"), "needs to specify a specific kind"))
		}
		if objref.APIVersion == "" && needsResourceType {
			errs = append(errs, field.Required(fld.Child("apiVersion"), "needs to specify a specific apiVersion"))
		}
		if objref.Name == "" {
			errs = append(errs, field.Required(fld.Child("name"), "needs to specify a resource name"))
		}
	}
	return errs
}

//ValidateCommonObject common validations for objects in katanomi,
// includes, name, annotations etc.
func ValidateCommonObject(obj metav1.Object) field.ErrorList {
	errs := field.ErrorList{}

	errs = append(errs, ValidateName(obj, field.NewPath("metadata", "name"))...)

	errs = append(errs, v1validation.ValidateLabels(obj.GetLabels(), field.NewPath("metadata", "labels"))...)

	errs = append(errs, apimachineryvalidation.ValidateAnnotations(obj.GetAnnotations(), field.NewPath("metadata", "annotations"))...)

	return errs
}

// ValidateItemName validates a name of an item in a slice. this is used in
//  Integration resources,  volumes and etc
func ValidateItemName(name string, prefix bool, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	errList := apimachineryvalidation.ValidateNamespaceName(name, prefix)
	if len(errList) > 0 {
		for _, errStr := range errList {
			errs = append(errs, field.Invalid(fld, name, errStr))
		}
	}
	return errs
}

// ReturnInvalidError returns a Invalid error if the error list is not empty
func ReturnInvalidError(gk schema.GroupKind, name string, errs field.ErrorList) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.NewInvalid(gk, name, errs)

}
