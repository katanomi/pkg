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

package url

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/apis"
)

var _ = Describe("Test.MatchGitURLPrefix", func() {

	DescribeTable("MatchGitURLPrefix",
		func(source, target string, expected bool) {
			s, err := apis.ParseURL(source)
			Expect(err).To(BeNil())
			t, err := apis.ParseURL(target)
			Expect(err).To(BeNil())
			actual := MatchGitURLPrefix(s, t)
			Expect(actual).Should(Equal(expected))
		},
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg.git", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "https://github.com/", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "http://github.com", true),
		Entry("test suffix .git", "https://github.com/katanomi/pkg.git", "github.com", false),
		//
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg/", false),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/pkg", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/katanomi/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg.git", "https://github.com/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg", "https://github.com/katanomi/", true),
		Entry("test suffix /", "https://github.com/katanomi/pkg", "http://github.com", true),
		//
		Entry("test suffix /", "https://github.com/katanomi/pkg", "https://github.com/katanomi/pkg.git", true),
	)
})
