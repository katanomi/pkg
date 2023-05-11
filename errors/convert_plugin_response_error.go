/*
Copyright 2023 The Katanomi Authors.

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
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ConvertToPluginK8StatusError converts a http response from plugin and an error into a kubernetes api error
// if the error is a kubernetes api error return it directly
// if the error is a common error then convert it to toolServiceUnavailable or Unauthorized error
// ctx is the basic context, response is the response object from tool sdk, err is the returned error
// names supports one optional name to be given and will be attributed as the resource name in the returned error
func ConvertToPluginK8StatusError(ctx context.Context, response *http.Response, err error, gvk schema.GroupVersionKind, names ...string) error {
	if err == nil {
		return err
	}
	if _, ok := err.(errors.APIStatus); ok {
		return err
	}

	var httpres *http.Response
	if response != nil {
		httpres = response
	}

	return commonErrorReason(ConvertResponseError(ctx, httpres, err, gvk, names...))
}

func commonErrorReason(err error) error {
	if statusErr, ok := err.(*errors.StatusError); ok {
		if statusErr.ErrStatus.Code == http.StatusUnauthorized {
			statusErr.ErrStatus.Reason = StatusReasonUnauthorized
		}
		if statusErr.ErrStatus.Code >= http.StatusInternalServerError {
			statusErr.ErrStatus.Reason = StatusReasonToolServiceUnavailable
		}
		return statusErr
	}
	return err
}
