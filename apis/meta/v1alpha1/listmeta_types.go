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

package v1alpha1

import (
	"regexp"

	"github.com/katanomi/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SortBy string

const (
	NameSortKey        SortBy = "name"
	UpdatedTimeSortKey SortBy = "updatedTime"
	CreatedTimeSortKey SortBy = "createdTime"

	// SearchValueKey The key of the keyword used for the search
	SearchValueKey = "searchValue"
)

type SortOrder string

const (
	OrderDesc SortOrder = "desc"
	OrderAsc  SortOrder = "asc"
)

// ListMeta extension of the metav1.ListMeta with paging related data
type ListMeta struct {
	metav1.ListMeta `json:",inline"`

	// TotalItems returned in the list
	TotalItems int `json:"totalItems"`

	// Current page number
	// +optional
	Page *int `json:"page,omitempty"`

	// Current items per page
	// +optional
	ItemsPerPage *int `json:"itemsPerPage,omitempty"`

	// Total number of pages
	// +optional
	TotalPages *int `json:"totalPages,omitempty"`
}

// ListOptions options for list
type ListOptions struct {
	Pager

	// All if true, return all items
	All bool `json:"all"`

	// Custom search options
	// +optional
	Search map[string][]string `json:",inline"`

	// Subresoures for listing
	// will only work for lists that support this feature
	// when not supported, this values will be ignored
	// +optional
	SubResources []string `json:"subResources"`

	// Sort for listing
	Sort []SortOptions `json:"sort"`
}

// GetSearchFirstElement get first element by key that in search map
func (opt *ListOptions) GetSearchFirstElement(key string) (value string) {
	if valueList, ok := opt.Search[key]; ok {
		if len(valueList) != 0 {
			value = valueList[0]
		}
	}
	return
}

// DefaultPager when ListOption page less than zero, set default value.
func (opt *ListOptions) DefaultPager() {
	if opt.ItemsPerPage < 1 {
		opt.ItemsPerPage = common.DefaultPerPage
	}

	if opt.Page < 1 {
		opt.Page = common.DefaultPage
	}
}

// MatchSearchValue match search value
func (opt *ListOptions) MatchSearchValue(name string) bool {
	values, ok := opt.Search[SearchValueKey]
	if !ok || len(values) == 0 {
		return true
	}
	for _, value := range values {
		if match, _ := regexp.MatchString(value, name); match {
			return true
		}
	}
	return false
}

// MatchSubResource match subresource
func (opt *ListOptions) MatchSubResource(name string) bool {
	if len(opt.SubResources) == 0 {
		return true
	}
	for _, sub := range opt.SubResources {
		if sub == name {
			return true
		}
	}
	return false
}

// SortOptions options for sort
type SortOptions struct {
	// SortBy field
	SortBy SortBy `json:"sortBy"`

	// Order sorted is 'asc' or 'desc'
	Order SortOrder `json:"order"`
}

// RepositoryOptions list repositroy path params
type RepositoryOptions struct {
	// project name
	Project string         `json:"project"`
	SubType ProjectSubType `json:"subType"`
}

// ArtifactOptions path params
type ArtifactOptions struct {
	RepositoryOptions

	// repository name
	Repository string `json:"repository"`

	// artifact name
	Artifact string `json:"artifact"`
}

// ArtifactTagOptions path params
type ArtifactTagOptions struct {
	ArtifactOptions

	// repository name
	Tag string `json:"tag"`
}

// IssueOptions path params
type IssueOptions struct {
	// Project identity name
	Identity string `json:"identity"`

	// Issue id
	IssueId string `json:"issueId"`

	// Issue branch
	Branch string `json:"branch"`
}

type UserOptions struct {
	// Project identity
	Project string `json:"project"`

	// Group identity
	Group string `json:"group"`

	// User indentity
	UserId string `json:"userId"`
}

// GetSearchValue get search value from option
// use `searchValue` instead of `name`
func GetSearchValue(option ListOptions) string {
	if value := option.GetSearchFirstElement(SearchValueKey); value != "" {
		return value
	}

	// Deprecated
	if value := option.GetSearchFirstElement("name"); value != "" {
		return value
	}

	return ""
}
