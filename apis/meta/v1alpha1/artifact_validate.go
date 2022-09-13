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
	"context"
	"net/url"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (p *ArtifactParameterSpec) Validate(ctx context.Context, path *field.Path) (errs field.ErrorList) {
	types := []string{
		string(OCIHelmChartArtifactParameterType),
		string(OCIContainerImageArtifactParameterType),
	}

	t := path.Child("type")
	if p.Type == "" {
		errs = append(errs, field.Required(t, validation.EmptyError()))
	} else {
		hasType := false
		for _, t := range types {
			if string(p.Type) == t {
				hasType = true
				break
			}
		}
		if !hasType {
			errs = append(errs, field.NotSupported(path.Child("type"), p.Type, types))
		}
	}

	uri := path.Child("URI")
	if p.URI == "" {
		errs = append(errs, field.Required(uri, validation.EmptyError()))
	} else {
		errs = append(errs, validateURI(p.URI, uri)...)
	}

	return
}

func validateURI(uri string, path *field.Path) (errs field.ErrorList) {
	if _, err := url.ParseRequestURI("http://" + uri); err != nil {
		errs = append(errs, field.Invalid(path, uri, err.Error()))
	}
	return
}
