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
	goerrors "errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//AsAPIError returns an error as a apimachinary api error
func AsAPIError(err error) error {
	reason := errors.ReasonForError(err)
	if reason == metav1.StatusReasonUnknown {
		err = errors.NewInternalError(err)
	}
	return err
}

// AsStatusCode returns the code from a errors.APIStatus, if not compatible will return InternalServerError
func AsStatusCode(err error) int {
	if status := errors.APIStatus(nil); goerrors.As(err, &status) {
		return int(status.Status().Code)
	}
	return http.StatusInternalServerError
}

// HandleError handles error in requests
func HandleError(req *restful.Request, resp *restful.Response, err error) {
	err = AsAPIError(err)
	status := AsStatusCode(err)
	resp.WriteError(status, err)
}
