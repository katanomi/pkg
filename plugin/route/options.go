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
	"time"

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
	FetchAllQueryKey     = "all"
)

// ParseTimeCursorFormRequest parse cursor from request
func ParseTimeCursorFormRequest(req *restful.Request) *metav1alpha1.TimeCursor {
	cursorStr := req.QueryParameter("continue")
	if cursorStr != "" {
		cursor, err := metav1alpha1.ParseTimeCursor(cursorStr)
		if err == nil {
			return cursor
		}
	}

	pager := ParsePagerFromRequest(req)
	tc := &metav1alpha1.TimeCursor{
		Pager:        pager,
		QueryStartAt: time.Now().Unix(),
	}
	return tc
}

// ParsePagerFromRequest parse the paging params from request
func ParsePagerFromRequest(req *restful.Request) metav1alpha1.Pager {
	p := metav1alpha1.Pager{}
	itemsPerPage := req.QueryParameter(ItemsPerPageQueryKey)
	if v, err := strconv.Atoi(itemsPerPage); err == nil {
		p.ItemsPerPage = v
	}
	page := req.QueryParameter(PageQueryKey)
	if v, err := strconv.Atoi(page); err == nil {
		p.Page = v
	}
	p.ItemsPerPage = p.GetPageLimit()
	p.Page = p.GetPage()
	return p
}

// GetListOptionsFromRequest returns ListOptions based on a request
func GetListOptionsFromRequest(req *restful.Request) (opts metav1alpha1.ListOptions) {
	pager := ParsePagerFromRequest(req)
	opts.Page = pager.Page
	opts.ItemsPerPage = pager.ItemsPerPage
	opts.Search = req.Request.URL.Query()
	if _, exist := opts.Search[FetchAllQueryKey]; exist {
		opts.All = true
	}
	delete(opts.Search, PageQueryKey)
	delete(opts.Search, ItemsPerPageQueryKey)
	delete(opts.Search, FetchAllQueryKey)

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
