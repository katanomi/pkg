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

package artifacts

import (
	"context"
	"fmt"
	"strings"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/regex"
)

const (
	// indicates pass filter result
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	PassFilter FilterResult = "pass"
	// indicates fail filter result
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	FailFilter FilterResult = "fail"
	// indicates no filter result
	//
	// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
	NoFilter FilterResult = "no_filter"
)

// FilterResult has the result of the filtering operation.
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type FilterResult string

// NewFilterResult will parse bool to filterresult
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewFilterResult(matched bool) FilterResult {
	if matched {
		return PassFilter
	}
	return FailFilter
}

// And will use AND logic between two filter result
// if any one is NoFilter will return other one
// if all are pass will return pass
// if one is fail will return fail
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (x FilterResult) And(y FilterResult) FilterResult {
	if x == NoFilter {
		return y
	}
	if y == NoFilter {
		return x
	}
	if x == PassFilter && y == PassFilter {
		return PassFilter
	}
	return FailFilter
}

// Or will use OR logic between two filter result
// if any one is NoFilter will return other one
// if any one is pass will return pass
// if all items are fail will return fail
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (x FilterResult) Or(y FilterResult) FilterResult {
	if x == NoFilter {
		return y
	}
	if y == NoFilter {
		return x
	}
	if x == PassFilter || y == PassFilter {
		return PassFilter
	}
	return FailFilter
}

// Not if current FilterResult is pass , will retrun fail
// if current FilterResult is faile, will return pass
// if current FilterResult is no filter , will return no filter
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (x FilterResult) Not() FilterResult {

	switch x {
	case PassFilter:
		return FailFilter
	case FailFilter:
		return PassFilter
	case NoFilter:
		return NoFilter
	default:
		return FailFilter
	}
}

// Filter is artifact filter interfact
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type Filter interface {
	// Filter compute the predicate on the provided event and returns the result of the matching
	// message will be match details when not pass
	Filter(ctx context.Context, artifact v1alpha1.Artifact) (res FilterResult, message string, err error)
}

// Filters is a wrapper that runs each filter and performs the and
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type Filters []Filter

// Filter compute the predicate on the provided event and returns the result of the matching
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (filters Filters) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	res := NoFilter
	for _, f := range filters {
		if f == nil {
			continue
		}

		itemRes, message, err := f.Filter(ctx, artifact)
		if err != nil {
			return FailFilter, "", err
		}
		res = res.And(itemRes)
		// Short circuit to optimize it
		if res == FailFilter {
			return FailFilter, message, nil
		}
	}

	return res, "", nil
}

var _ Filter = Filters{}

// NewFilter create Filter by filters, will use AND logic between items
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewFilter(filters ...v1alpha1.ArtifactFilterSet) Filter {

	res := Filters{}

	for _, filter := range filters {
		if len(filter.All) > 0 {
			res = append(res, NewAllFilter(filter.All))
		}

		if len(filter.Any) > 0 {
			res = append(res, NewAnyFilter(filter.Any))
		}
	}

	return res
}

// NewAllFilter will construct all filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewAllFilter(all []v1alpha1.ArtifactFilter) Filter {
	return AllFilter{items: all}
}

// AllFilter will use AND logic between items
// in other words, it equals item[0] && item[1] && ... item[n] (n>=0)
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type AllFilter struct {
	items []v1alpha1.ArtifactFilter
}

// Filter will filter artifacts by all filters
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (f AllFilter) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	if len(f.items) == 0 {
		return NoFilter, "", nil
	}

	for _, item := range f.items {
		filter := NewArtifactFilter(item)
		if filter == nil {
			continue
		}

		v, message, err := filter.Filter(ctx, artifact)
		if err != nil {
			return FailFilter, "", err
		}

		if v == FailFilter {
			return FailFilter, message, nil
		}
	}

	return PassFilter, "", nil
}

// NewAnyFilter will construct filters using OR logic
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewAnyFilter(any []v1alpha1.ArtifactFilter) Filter {
	return AnyFilter{items: any}
}

// AnyFilter will use OR logic between items
// in other words, it equals item[0] || item[1] || ... item[n] (n>=0)
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type AnyFilter struct {
	items []v1alpha1.ArtifactFilter
}

// Filter will filter artifacts by filters using OR logic
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (f AnyFilter) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	if len(f.items) == 0 {
		return NoFilter, "", nil
	}

	var messages = []string{}
	for _, item := range f.items {
		filter := NewArtifactFilter(item)
		if filter == nil {
			continue
		}

		v, message, err := filter.Filter(ctx, artifact)
		if err != nil {
			return FailFilter, "", err
		}
		if v == PassFilter {
			return PassFilter, "", nil
		}
		messages = append(messages, message)
	}

	return FailFilter, strings.Join(messages, ","), nil
}

// NewArtifactFilter will construct a ArtifactFilter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewArtifactFilter(filter v1alpha1.ArtifactFilter) Filters {
	filters := Filters{}

	if len(filter.Labels) > 0 {
		for _, item := range filter.Labels {
			filters = append(filters, NewLabelFilter(item))
		}
	}
	if len(filter.Tags) > 0 {
		for _, item := range filter.Tags {
			filters = append(filters, NewTagFilter(item))
		}
	}
	if len(filter.Envs) > 0 {
		for _, item := range filter.Envs {
			filters = append(filters, NewEnvFilter(item))
		}
	}

	return filters
}

// LabelFilter artifact label filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type LabelFilter struct {
	v1alpha1.ArtifactLabelFilter
}

// NewLabelFilter will construct artifact label filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewLabelFilter(f v1alpha1.ArtifactLabelFilter) Filter {
	return &LabelFilter{f}
}

// Filter will filter artifacts by artifact labels
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (filter LabelFilter) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	if len(filter.Regex) == 0 {
		return NoFilter, "", nil
	}

	p, err := artifact.ParseProperties()
	if err != nil {
		return FailFilter, "", err
	}

	if len(p.ExtraAttrs.Config.Labels) == 0 {
		return FailFilter, "labels are empty", nil
	}
	val, ok := p.ExtraAttrs.Config.Labels[filter.Name]
	if !ok {
		return FailFilter, fmt.Sprintf("not contains label name: %s", filter.Name), nil
	}

	matched, err := regex.Regexes(filter.Regex).MatchString(val)
	if err != nil {
		return FailFilter, "", err
	}

	message := ""
	if !matched {
		message = fmt.Sprintf(" '%s' value in labels are not matched '%v' ", filter.Name, filter.Regex)
	}
	return NewFilterResult(matched), message, nil
}

var _ Filter = EnvFilter{}

// EnvFilter represents artifact env filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type EnvFilter struct {
	v1alpha1.ArtifactEnvFilter
}

// NewEnvFilter will construct artifact env filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewEnvFilter(f v1alpha1.ArtifactEnvFilter) Filter {
	return &EnvFilter{f}
}

// Filter will filter artifacts by artifact envs
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (filter EnvFilter) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	p, err := artifact.ParseProperties()
	if err != nil {
		return NoFilter, "", err
	}

	if len(filter.Regex) == 0 {
		return NoFilter, "", nil
	}

	envs := v1alpha1.ParseEnvs(p.ExtraAttrs.Config.Env)
	if len(envs) == 0 {
		return FailFilter, "envs are empty", nil
	}

	val, ok := envs[filter.Name]
	if !ok {
		return FailFilter, fmt.Sprintf("not contains env name: %s", filter.Name), nil
	}

	matched, err := regex.Regexes(filter.Regex).MatchString(val)
	if err != nil {
		return FailFilter, "", err
	}

	message := ""
	if !matched {
		message = fmt.Sprintf(" '%s' value in envs are not matched '%v' ", filter.Name, filter.Regex)
	}

	return NewFilterResult(matched), message, nil
}

var _ Filter = TagFilter{}

// TagFilter represents artifact tag filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
type TagFilter struct {
	v1alpha1.ArtifactTagFilter
}

// NewTagFilter will construct artifact tag filter
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func NewTagFilter(f v1alpha1.ArtifactTagFilter) Filter {
	return &TagFilter{f}
}

// Filter will filter artifacts by artifact tags
//
// Deprecated: use pkg/apis/artifacts/v1alpha1 instead
func (filter TagFilter) Filter(ctx context.Context, artifact v1alpha1.Artifact) (FilterResult, string, error) {
	p, err := artifact.ParseProperties()
	if err != nil {
		return NoFilter, "", err
	}

	if len(filter.Regex) == 0 {
		return NoFilter, "", nil
	}

	if len(p.Tags) == 0 {
		return FailFilter, "tags is empty", nil
	}

	tags := []string{}
	for _, tag := range p.Tags {
		tags = append(tags, tag.Name)
	}

	matchedStrs, err := regex.Regexes(filter.Regex).MatchAnyString(tags...)
	if err != nil {
		return NoFilter, "", err
	}

	if len(matchedStrs) == 0 {
		return FailFilter, fmt.Sprintf(" tags are not matched '%v' ", filter.Regex), nil
	}

	return PassFilter, "", nil
}
