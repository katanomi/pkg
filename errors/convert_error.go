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

// Package errors contains useful functionality for conversion errors
package errors

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"knative.dev/pkg/apis"
)

func ConvertKnativeFieldErrorToInternalError(err *apis.FieldError) *field.Error {
	if err == nil {
		return nil
	}
	return field.InternalError(field.NewPath(""), err)
}

func ConvertKnativeFieldErrorToErrorList(err *apis.FieldError) field.ErrorList {
	if err == nil {
		return nil
	}
	return field.ErrorList{
		ConvertKnativeFieldErrorToInternalError(err),
	}
}
