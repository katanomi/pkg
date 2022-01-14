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
	"github.com/go-resty/resty/v2"
	kerrors "github.com/katanomi/pkg/errors"
)

// GetErrorFromResponse returns an error based on the response. Will do the best effort to convert
// error responses into apimachinery errors
func GetErrorFromResponse(resp *resty.Response, err error) error {
	if resp != nil && resp.IsError() {
		return kerrors.AsStatusError(resp)
	}
	return err
}
