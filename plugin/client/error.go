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

package client

import (
	"net/http"

	"github.com/katanomi/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsNotImplementedError returns true if the plugin not implement the specified interface
func IsNotImplementedError(err error) bool {
	return errors.AsStatusCode(err) == http.StatusNotFound
}

// ResponseStatusErr is an error with `Status` type,
// used to handle plugin response
type ResponseStatusErr struct {
	metav1.Status
}

// Error implements error interface, output the error message
func (p ResponseStatusErr) Error() string {
	return p.Status.Message
}
