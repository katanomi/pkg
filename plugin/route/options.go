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

package route

import (
	"strconv"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// GetListOptionsFromRequest returns ListOptions based on a request
func GetListOptionsFromRequest(req *restful.Request) (opts metav1alpha1.ListOptions) {
	itemsPerPage := req.QueryParameter("itemsPerPage")
	if v, err := strconv.Atoi(itemsPerPage); err == nil {
		opts.ItemsPerPage = v
	}
	page := req.QueryParameter("page")
	if v, err := strconv.Atoi(page); err == nil {
		opts.Page = v
	}

	opts.Search = req.Request.URL.Query()
	delete(opts.Search, "page")
	delete(opts.Search, "itemsPerPage")
	return
}

// ListOptionsDocs adds list options query parameters to the documentation
func ListOptionsDocs(bldr *restful.RouteBuilder) *restful.RouteBuilder {
	// TODO: adds parameters to lists here
	return bldr
}

func GetPathParamsFromRequest(req *restful.Request, names ...string) (params metav1alpha1.PathParams) {
	params = make(metav1alpha1.PathParams)
	for _, name := range names {
		params[name] = req.PathParameter(name)
	}

	return
}
