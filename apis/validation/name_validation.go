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
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const dns1123UnderscoreLabelFmt string = "[A-Za-z0-9]([-A-Za-z0-9_/]*[A-Za-z0-9])?"
const dns1123UnderscoreLabelErrMsg string = "a lowercase RFC 1123 with underscore label must consist of lower case alphanumeric characters, '_' or '-', and must start and end with an alphanumeric character"

// DNS1123UnderscoreLabelMaxLength is a label's max length in DNS (RFC 1123)
const DNS1123UnderscoreLabelMaxLength int = validation.DNS1123LabelMaxLength

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123UnderscoreLabelFmt + "$")

const genericResourceNameFmt string = "[-./A-Za-z0-9_]+"
const genericResourceNameErrMsg string = "a resource name must consist of lower case alphanumeric characters, '_' or '-' or '/' or '.' "

var genericResourceNameRegexp = regexp.MustCompile("^" + genericResourceNameFmt + "$")

const resourceNameWithChineseFmt = "[-./\\sA-Za-z0-9_\\p{Han}]+"
const resourceNameWithChineseErrMsg string = "a resource name must consist of lower case alphanumeric characters, Chinese, '_' or '-' or '/' or '.' "

var resourceNameWithChineseRegexp = regexp.MustCompile("^[^\\s]" + resourceNameWithChineseFmt + "$")

// IsDNS1123UnderscoreLabel tests for a string that conforms to the definition of a label in
// DNS (RFC 1123) but accepting underscores in the middle.
func IsDNS1123UnderscoreLabel(value string) []string {
	var errs []string
	if len(value) > DNS1123UnderscoreLabelMaxLength {
		errs = append(errs, validation.MaxLenError(DNS1123UnderscoreLabelMaxLength))
	}
	if !dns1123LabelRegexp.MatchString(value) {
		errs = append(errs, validation.RegexError(dns1123UnderscoreLabelErrMsg, dns1123UnderscoreLabelFmt, "my-name", "123_abc"))
	}
	return errs
}

func IsGenericResourceName(value string) []string {
	var errs []string

	if !genericResourceNameRegexp.MatchString(value) {
		errs = append(errs, validation.RegexError(genericResourceNameErrMsg, genericResourceNameFmt, "my-name", "abc/123"))
	}
	return errs
}

// ValidateItemNameUnderscore validates a name of an item in a slice. this is used in
// resources,  volumes and etc
func ValidateItemNameUnderscore(name string, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	errList := IsDNS1123UnderscoreLabel(name)
	if len(errList) > 0 {
		for _, errStr := range errList {
			errs = append(errs, field.Invalid(fld, name, errStr))
		}
	}
	return errs
}

func ValidateGenericResourceName(name string, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	errList := IsGenericResourceName(name)
	if len(errList) > 0 {
		for _, errStr := range errList {
			errs = append(errs, field.Invalid(fld, name, errStr))
		}
	}
	return errs
}

func ValidateResourceNameWithChinese(name string, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	if !resourceNameWithChineseRegexp.MatchString(name) {
		err := validation.RegexError(resourceNameWithChineseErrMsg, resourceNameWithChineseFmt, "my-name", "abc/123/示例")
		errs = append(errs, field.Invalid(fld, name, err))
	}

	return errs
}
