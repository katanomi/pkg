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

package client

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/go-resty/resty/v2"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// SecretOpts provides a secret to be assigned to the request in the header
func SecretOpts(secret corev1.Secret) OptionFunc {
	return func(request *resty.Request) {
		auth := FromSecret(secret)

		request.SetHeader(PluginAuthHeader, string(auth.Type))
		dataBytes, _ := json.Marshal(auth.Secret)
		request.SetHeader(PluginSecretHeader, base64.StdEncoding.EncodeToString(dataBytes))
	}
}

// MetaOpts provides metadata for the request
func MetaOpts(meta Meta) OptionFunc {
	return func(request *resty.Request) {
		dataBytes, _ := json.Marshal(meta)
		request.SetHeader(PluginMetaHeader, base64.StdEncoding.EncodeToString(dataBytes))
	}
}

// ListOpts options for lists
func ListOpts(opts metav1alpha1.ListOptions) OptionFunc {
	return func(request *resty.Request) {
		if len(opts.Search) > 0 {
			for k, v := range opts.Search {
				for _, val := range v {
					request.SetQueryParam(k, val)
				}
			}
		}
		request.SetQueryParam("page", strconv.Itoa(opts.Page))
		request.SetQueryParam("itemsPerPage", strconv.Itoa(opts.ItemsPerPage))
	}
}

// QueryOpts query parameters for the request
func QueryOpts(params map[string]string) OptionFunc {
	return func(request *resty.Request) {
		request.SetQueryParams(params)
	}
}

// BodyOpts request body
func BodyOpts(body interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetBody(body)
	}
}

// ResultOpts request result automatically marshalled into object
func ResultOpts(dest interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetResult(dest)
	}
}

// ErrorOpts error response object
func ErrorOpts(err interface{}) OptionFunc {
	return func(request *resty.Request) {
		request.SetError(err)
	}
}

// HeaderOpts sets a header
func HeaderOpts(key, value string) OptionFunc {
	return func(request *resty.Request) {
		request.SetHeader(key, value)
	}
}
