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
	"context"
	"net/http"


	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ConvertResponseError converts a http response and an error into a kubernetes api error
// ctx is the basic context, response is the response object from gitlab sdk, err is the returned error
// gvk is the GroupVersionKind object with type meta for the object
// names supports one optional name to be given and will be attributed as the resource name in the returned error
func ConvertResponseError(ctx context.Context, response *http.Response, err error, gvk schema.GroupVersionKind, names ...string) error {
	if err == nil {
		return err
	}
	statusCode := http.StatusInternalServerError
	method := http.MethodGet
	name := ""
	if response != nil {
		statusCode = response.StatusCode
		method = response.Request.Method
	}
	if len(names) > 0 {
		name = names[0]
	} else if response != nil && response.Request != nil && response.Request.URL != nil {
		name = response.Request.URL.String()
	} else {
		// use default
	}

	return errors.NewGenericServerResponse(
		statusCode,
		method,
		schema.GroupResource{Group: gvk.Group, Resource: gvk.Kind},
		name, err.Error(),
		0,
		true)
}

