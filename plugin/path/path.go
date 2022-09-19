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

package path

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"net/url"
	"strings"
)

// Escape escapes the string, so it can be safely placed inside a URL path segment
//
// url.PathEscape will encode "/" to "%2F", but go-restful will decode it to "/" and return 405 error.
// so should replace special characters (including /) with %XX sequences first.
func Escape(path string) string {
	path = strings.Replace(path, "/", "%2F", 1)
	return url.PathEscape(path)
}

// Format formats the path with the given parameters
// escape the path parameters automatically
func Format(tmpl string, params ...string) string {
	list := make([]interface{}, 0, len(params))
	for _, p := range params {
		list = append(list, Escape(p))
	}
	return fmt.Sprintf(tmpl, list...)
}

// Parameter gets the path parameter from the request
func Parameter(request *restful.Request, key string) string {
	path := request.PathParameter(key)
	return strings.Replace(path, "%2F", "/", 1)
}
