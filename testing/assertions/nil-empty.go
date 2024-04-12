/*
Copyright 2024 The Katanomi Authors.

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

// Package assertions contains all the assertions that can be used in tests
// for general testing. All methods are Gomega compatible and can be used
// directly with gomega library
package assertions

import (
	"reflect"

	"github.com/onsi/gomega/gcustom"
	gtypes "github.com/onsi/gomega/types"
)

// BeNilOrEmpty generates a GomegaMatcher that checks if the value is nil or empty.
// works for pointers and to empty values
func BeNilOrEmpty() gtypes.GomegaMatcher {
	return gcustom.MakeMatcher(func(value any) (bool, error) {
		typeOf := reflect.TypeOf(value)
		if value == nil || typeOf == nil {
			return true, nil
		}

		valueOf := reflect.ValueOf(value)
		if valueOf.Kind() == reflect.Ptr && valueOf.IsNil() {
			return true, nil
		} else if valueOf.Kind() == reflect.Ptr {
			valueOf = valueOf.Elem()
			value = valueOf.Interface()
		}

		zeroValue := reflect.Zero(valueOf.Type()).Interface()
		return reflect.DeepEqual(zeroValue, value), nil
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} be nil or empty")
}
