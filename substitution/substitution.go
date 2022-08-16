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

// Package substitution contains useful functionality for validating variables
package substitution

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"knative.dev/pkg/apis"
)

const parameterSubstitutionSample = `[_a-zA-Z][_a-zA-Z0-9.-]*(\[\*\])?`

const BraceMatchingRegex = "(\\$(\\(%s\\.(?P<var>%s)\\)))"

type ValidateSingleVariableFunc func(vars sets.String, path string, name string, value string, locationName string, paramName string) *apis.FieldError

func TektonValidateSingleVariable(vars sets.String, path string, name, value, locationName, v string) *apis.FieldError {
	v = strings.TrimSuffix(v, "[*]")
	if !vars.Has(v) {
		return &apis.FieldError{
			Message: fmt.Sprintf("non-existent variable in %q for %s %s", value, locationName, name),
			Paths:   []string{path + "." + name},
		}
	}
	return nil
}

func VersionValidateSingleVariable(vars sets.String, path string, name, value, locationName, v string) *apis.FieldError {
	if !vars.Has(v) {
		return &apis.FieldError{
			Message: fmt.Sprintf("non-existent variable in %q for %s %s", value, locationName, v),
			Paths:   []string{path},
		}
	}
	return nil
}

func ValidateVariable(name, value, prefix, locationName, path string, vars sets.String,
	paramSubstitution string, validateFunc ValidateSingleVariableFunc) *apis.FieldError {
	if vs, present := extractVariablesFromString(value, prefix, paramSubstitution); present {
		for _, v := range vs {
			err := validateFunc(vars, path, name, value, locationName, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func extractVariablesFromString(s, prefix, paramSubstitution string) ([]string, bool) {
	pattern := fmt.Sprintf(BraceMatchingRegex, prefix, paramSubstitution)
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return []string{}, false
	}
	vars := make([]string, len(matches))
	for i, match := range matches {
		groups := matchGroups(match, re)
		// foo -> foo
		// foo.bar -> foo
		// foo.bar.baz -> foo

		// ---- this is tekton's implementation
		// we change it here because we want to validate
		// the sub variables, like foo.bar.baz
		// vars[i] = strings.SplitN(groups["var"], ".", 2)[0]

		// ---- takes the whole var as accepted type
		// allowing to validate foo.bar.baz variables
		vars[i] = groups["var"]
	}
	return vars, true
}

func matchGroups(matches []string, pattern *regexp.Regexp) map[string]string {
	groups := make(map[string]string)
	for i, name := range pattern.SubexpNames()[1:] {
		groups[name] = matches[i+1]
	}
	return groups
}

// GetValWithKeyFunc standard function to extract value maps from objects
type GetValWithKeyFunc func(ctx context.Context, path *field.Path) map[string]string
