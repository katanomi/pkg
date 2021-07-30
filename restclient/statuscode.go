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

package restclient

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// RESTClientGroupResource fake GroupResource to use errors api
var RESTClientGroupResource = schema.GroupResource{Group: "katanomi.dev", Resource: "RESTfulClient"}

// GetErrorFromResponse returns an error based on the response. Will do the best effort to convert
// error responses into apimachinery errors
func GetErrorFromResponse(resp *resty.Response, err error) error {
	if resp.IsError() {
		if err == nil {
			err = fmt.Errorf(resp.String())
		}
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return errors.NewBadRequest(err.Error())
		case http.StatusUnauthorized:
			return errors.NewUnauthorized(err.Error())
		case http.StatusMethodNotAllowed:
			return errors.NewMethodNotSupported(RESTClientGroupResource, "")
		case http.StatusInternalServerError:
			return errors.NewInternalError(err)
		case http.StatusRequestTimeout:
			return errors.NewTimeoutError(err.Error(), 0)
		default:
			return err
		}
	}
	return nil
}
