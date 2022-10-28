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

package validators

import (
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.URL.Validate", func() {
	var (
		urlValidator *URL
		urls         []string
		rootPath     = field.NewPath("root")
	)

	BeforeEach(func() {
		urlValidator = NewURL()
		urls = make([]string, 0)
	})

	Context("validate invalid url", func() {
		JustBeforeEach(func() {
			urls = []string{
				"cache_object://foo",
				"htt:// abc",
			}
		})
		It("you should return the same number of errors as the url", func() {
			errs := urlValidator.Validate(rootPath, urls...)
			Expect(len(errs)).To(Equal(len(urls)))
		})
	})

	Context("validate valid url", func() {
		JustBeforeEach(func() {
			urls = []string{
				"aaa",
				"aaa/bbb",
				"aaa/bbb:latest",
				"aaa/bbb:v0.0.1",
				"http://aaa/bbb:v0.0.1",
				"https://aaa/bbb:v0.0.1",
			}
		})
		It("you should return the same number of errors as the url", func() {
			errs := urlValidator.Validate(rootPath, urls...)
			Expect(len(errs)).To(Equal(0))
		})
	})

	Context("customize error message", func() {
		It("got the customize error message", func() {
			urlValidator.SetErrMsg("customized error message")
			errs := urlValidator.Validate(rootPath, "htt:// abc")
			Expect(errs.ToAggregate().Error()).To(ContainSubstring("customized error message"))
		})
	})

	Context("customize validate function", func() {
		It("got the customize error message", func() {
			urlValidator.SetValidate(func(url *url.URL) (ok bool, errMsg string) {
				if url.Host == "" {
					return false, "host is empty"
				}
				return true, ""
			})
			errs := urlValidator.Validate(rootPath, "abc")
			Expect(errs.ToAggregate().Error()).To(ContainSubstring("host is empty"))
		})
	})
})
