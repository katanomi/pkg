/*
Copyright 2023 The Katanomi Authors.

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
	. "github.com/onsi/gomega"
	"testing"
)

func TestFindMatchedTagByPrefix(t *testing.T) {
	ctx := context.TODO()
	table := map[string]struct {
		OldTag      string
		NewTags     []string
		Options     FindMatchedTagByPrefixOptions
		ExpectedTag string
	}{
		"has multiple matching prefixes": {
			"v1.8.0-alpha.123456",
			[]string{
				"v1.9.0-alpha.12",
				"v1.8.1-alpha.12345",
				"v1.8.0-beta.1234",
				"v1.8.0-alpha.654",
			},
			FindMatchedTagByPrefixOptions{},
			"v1.8.0-alpha.654",
		},
		"matches next patch version": {
			"v1.8.0-alpha.34",
			[]string{
				"v1.9.0-beta.12",
				"v1.8.1-delta.654",
				"v1.8.1-beta.34",
				"v1.9.1-beta.1234",
			},
			FindMatchedTagByPrefixOptions{5},
			"v1.8.1-delta.654",
		},
		"no matches": {
			"v1.8.0-alpha.123456",
			[]string{
				"alpine-3.4-abc",
				"ubuntu-3214.abc",
				"debian-xd123",
				"centos-x123",
			},
			FindMatchedTagByPrefixOptions{3},
			"",
		},
		"min characters not matching": {
			"v1.8.0-alpha.123456",
			[]string{
				"v1.8.1-alpha.123456",
				"v1.8.2-alpha.123456",
				"v1.8.3-alpha.123456",
				"v1.8.4-alpha.123456",
			},
			FindMatchedTagByPrefixOptions{6},
			"",
		},
	}
	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			tag, _ := FindMatchedTagByPrefix(ctx, test.OldTag, test.Options, test.NewTags...)
			g.Expect(tag).To(Equal(test.ExpectedTag))
		})
	}
}
