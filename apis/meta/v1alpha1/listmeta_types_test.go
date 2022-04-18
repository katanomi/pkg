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

package v1alpha1

import (
	"testing"

	"github.com/katanomi/pkg/common"
	. "github.com/onsi/gomega"
)

func TestGetSearchValue(t *testing.T) {
	opt := ListOptions{
		Search: map[string][]string{SearchValueKey: {"test", "test2"}},
	}
	value := opt.GetSearchFirstElement(SearchValueKey)
	if value != "test" {
		t.Fail()
		return
	}
	opt.Search = map[string][]string{SearchValueKey: {}}
	value = opt.GetSearchFirstElement(SearchValueKey)
	if value != "" {
		t.Fail()
		return
	}
	opt.Search = map[string][]string{}
	value = opt.GetSearchFirstElement(SearchValueKey)
	if value != "" {
		t.Fail()
		return
	}
}

func TestListOptions_DefaultPager(t *testing.T) {
	opt := ListOptions{
		Page:         0,
		ItemsPerPage: 0,
	}

	opt.DefaultPager()
	Equal(opt.Page).Match(common.DefaultPage)
	Equal(opt.ItemsPerPage).Match(common.DefaultPerPage)

	opt1 := ListOptions{
		Page:         10,
		ItemsPerPage: 0,
	}
	opt1.DefaultPager()
	opt1.DefaultPager()

	Equal(opt1.Page).Match(10)
	Equal(opt1.ItemsPerPage).Match(common.DefaultPerPage)

	opt2 := ListOptions{
		Page:         10,
		ItemsPerPage: 30,
	}
	opt2.DefaultPager()

	Equal(opt2.Page).Match(10)
	Equal(opt2.ItemsPerPage).Match(30)
}
