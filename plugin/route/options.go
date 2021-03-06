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
	"context"
	"net/url"
	"strconv"
	"strings"

	"knative.dev/pkg/logging"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	pclient "github.com/katanomi/pkg/plugin/client"
)

const (
	PageQueryKey         = "page"
	ItemsPerPageQueryKey = "itemsPerPage"
	SortQueryKey         = "sortBy"
	SinceQueryKey        = "since"
	UntilQueryKey        = "until"
)

// GetListOptionsFromRequest returns ListOptions based on a request
func GetListOptionsFromRequest(req *restful.Request) (opts metav1alpha1.ListOptions) {
	itemsPerPage := req.QueryParameter(ItemsPerPageQueryKey)
	if v, err := strconv.Atoi(itemsPerPage); err == nil {
		opts.ItemsPerPage = v
	}
	page := req.QueryParameter(PageQueryKey)
	if v, err := strconv.Atoi(page); err == nil {
		opts.Page = v
	}

	opts.Search = req.Request.URL.Query()
	delete(opts.Search, PageQueryKey)
	delete(opts.Search, ItemsPerPageQueryKey)

	subResourcesHeader := req.HeaderParameter(pclient.PluginSubresourcesHeader)
	if strings.TrimSpace(subResourcesHeader) != "" {
		opts.SubResources = strings.Split(subResourcesHeader, ",")
	}

	sortBy := req.QueryParameter(SortQueryKey)
	if sortBy == "" {
		return
	}
	opts.Sort = HandleSortQuery(req.Request.Context(), sortBy)
	return
}

// HandleSortQuery parse the sorting parameters in the query and return a list of sorting parameters
func HandleSortQuery(ctx context.Context, sortBy string) []metav1alpha1.SortOptions {
	logger := logging.FromContext(ctx).Named("SortParamHandler")
	sortList := []metav1alpha1.SortOptions{}
	sortInfoList := strings.Split(sortBy, ",")
	if len(sortInfoList)%2 != 0 {
		logger.Errorw("sortBy is expected to be in pairs ignoring sort by...", SortQueryKey, url.QueryEscape(sortBy), "sortInfoListLen", len(sortInfoList))
		return sortList
	}
	for i, v := range sortInfoList {
		if i%2 == 0 {
			switch metav1alpha1.SortOrder(v) {
			case metav1alpha1.OrderDesc, metav1alpha1.OrderAsc:
				sortList = append(sortList, metav1alpha1.SortOptions{Order: metav1alpha1.SortOrder(v)})
			default:
				sortList = []metav1alpha1.SortOptions{}
				logger.Errorw("unknown order type", SortQueryKey, url.QueryEscape(sortBy))
				return sortList
			}
		} else {
			sortList[len(sortList)-1].SortBy = metav1alpha1.SortBy(v)
		}
	}
	return sortList
}

// ListOptionsDocs adds list options query parameters to the documentation
func ListOptionsDocs(bldr *restful.RouteBuilder) *restful.RouteBuilder {
	// TODO: adds parameters to lists here
	return bldr.Param(restful.QueryParameter(ItemsPerPageQueryKey, "items to be returned in a page")).
		Param(restful.QueryParameter(PageQueryKey, "page to be returned"))
}
