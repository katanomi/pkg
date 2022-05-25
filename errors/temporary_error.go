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
	"k8s.io/apimachinery/pkg/api/errors"
)

// IsTemporaryError this will check against a list of know
// non-temporary related errors like 403-Forbidden etc
// all others will be considered to be temporary
func IsTemporaryError(err error) bool {
	if tempErr, ok := err.(TemporaryError); ok {
		return tempErr.Temporary()
	}
	// fallback to verify other errors such as apimachinery errors
	switch {
	case errors.IsForbidden(err),
		errors.IsInvalid(err),
		errors.IsAlreadyExists(err),
		errors.IsNotFound(err),
		errors.IsMethodNotSupported(err),
		errors.IsGone(err),
		errors.IsRequestEntityTooLargeError(err),
		errors.IsNotAcceptable(err),
		errors.IsUnauthorized(err):
		// TODO: this will not be possible for now, should fail
		return false
	default:
		return true
	}
}

// TemporaryError is an error and has a Temporary function
type TemporaryError interface {
	error
	Temporary() bool
}
