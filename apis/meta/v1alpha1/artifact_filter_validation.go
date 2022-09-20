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

package v1alpha1

import (
	"regexp"

	"github.com/katanomi/pkg/apis/validation"
	kvalidation "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (a ArtifactFilterRegexList) Validate(path *field.Path) (errs field.ErrorList) {
	if len(a) == 0 {
		errs = append(errs, field.Required(path, kvalidation.EmptyError()))
	}

	for i, regex := range a {
		if regex == "" {
			errs = append(errs, field.Required(path.Index(i), kvalidation.EmptyError()))
		} else if _, err := regexp.Compile(regex); err != nil {
			errs = append(errs, field.Invalid(path.Index(i), regex, err.Error()))
		}
	}

	return
}

func (a *ArtifactTagFilter) Validate(path *field.Path) (errs field.ErrorList) {
	errs = append(errs, a.Regex.Validate(path.Child("regex"))...)

	return
}

func (a *ArtifactEnvFilter) Validate(path *field.Path) (errs field.ErrorList) {
	errs = append(errs, validateName(a.Name, path.Child("name"))...)
	errs = append(errs, a.Regex.Validate(path.Child("regex"))...)

	return
}

func (a *ArtifactLabelFilter) Validate(path *field.Path) (errs field.ErrorList) {
	errs = append(errs, validateName(a.Name, path.Child("name"))...)
	errs = append(errs, a.Regex.Validate(path.Child("regex"))...)

	return
}

func (a *ArtifactFilter) Validate(path *field.Path) (errs field.ErrorList) {
	for i, filter := range a.Envs {
		errs = append(errs, filter.Validate(path.Child("envs").Index(i))...)
	}

	for i, filter := range a.Tags {
		errs = append(errs, filter.Validate(path.Child("tags").Index(i))...)
	}

	for i, filter := range a.Labels {
		errs = append(errs, filter.Validate(path.Child("labels").Index(i))...)
	}

	return
}

func (a *ArtifactFilterSet) Validate(path *field.Path) (errs field.ErrorList) {
	if len(a.Any) > 0 && len(a.All) > 0 {
		errs = append(errs, field.Invalid(path, "", "only one of the two fields [all, any] is required"))
	}

	for i, filter := range a.Any {
		errs = append(errs, filter.Validate(path.Child("any").Index(i))...)
	}

	for i, filter := range a.All {
		errs = append(errs, filter.Validate(path.Child("all").Index(i))...)
	}

	return
}

func validateName(name string, path *field.Path) (errs field.ErrorList) {
	if name == "" {
		errs = append(errs, field.Required(path, kvalidation.EmptyError()))
	} else {
		errs = append(errs, validation.ValidateGenericResourceName(name, path)...)
	}

	return
}
