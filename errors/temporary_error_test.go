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

package errors

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type someTemporaryError struct {
	temporary bool
}

func (s someTemporaryError) Temporary() bool {
	return s.temporary
}

func (s someTemporaryError) Error() string {
	return fmt.Sprintf("temporary error? %t", s.temporary)
}

func TestIsTemporaryError(t *testing.T) {

	table := map[string]struct {
		Error    error
		Expected bool
	}{
		"Forbidden is not temporary": {
			errors.NewForbidden(schema.GroupResource{}, "error", fmt.Errorf("some error")),
			false,
		},
		"Not found is not temporary": {
			errors.NewNotFound(schema.GroupResource{}, "error"),
			false,
		},
		"Unauthorized is not temporary ": {
			errors.NewUnauthorized("error"),
			false,
		},
		"BadRequest is not temporary ": {
			errors.NewBadRequest("bad request"),
			false,
		},
		// ... and so on

		"Internal server error is a temporary error": {
			errors.NewInternalError(fmt.Errorf("some error")),
			true,
		},
		"Some random error is temporary": {
			fmt.Errorf("random error"),
			true,
		},
		"temporary error true": {
			someTemporaryError{temporary: true},
			true,
		},
		"temporary error false": {
			someTemporaryError{temporary: false},
			false,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {

			result := IsTemporaryError(test.Error)
			if result != test.Expected {
				t.Errorf("result is unexpected %t != %t", test.Expected, result)
			}
		})
	}
}
