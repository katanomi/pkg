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

package v1alpha1

import (
	kvalidation "github.com/katanomi/pkg/apis/validation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"knative.dev/pkg/apis"
)

// ResourceURI stores a resource URI together with secret references
// for usage
type ResourceURI struct {
	// Unique resource identification for webhook resource, i.e
	// a code repository uri: https://github.com/katanomi/core
	URI *apis.URL `json:"uri"`

	// SecretRef stores a reference to a secret object
	// that contain authentication data for the described resource
	SecretRef *corev1.ObjectReference `json:"secretRef"`
}

// Validate basic validation for ResourceURI
func (r *ResourceURI) Validate(path *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	errs = append(errs, kvalidation.ValidateURL(r.URI, path.Child("uri"))...)
	errs = append(errs, kvalidation.ValidateObjectReference(r.SecretRef, false, false, path.Child("secretRef"))...)
	return errs
}
