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

package secret

import (
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestResourcePathFormat_FormatPathByScene", func() {
	var (
		pathFmtJson    string
		subPathFmtJson string
		scope          string
		scene          metav1alpha1.ResourcePathScene
		newScope       string

		r *ResourcePathFormat
	)

	BeforeEach(func() {
		pathFmtJson = ""
		subPathFmtJson = ""
		scope = ""
		newScope = ""
		scene = ""
		r = &ResourcePathFormat{}
	})

	JustBeforeEach(func() {
		r = NewResourcePathFormat(pathFmtJson, subPathFmtJson)
		newScope = r.FormatPathByScene(scene, scope)
	})
	Context("json string is invalid", func() {
		BeforeEach(func() {
			pathFmtJson = "invalid json"
			subPathFmtJson = "invalid json"
		})
		Context("scope is a sub resource", func() {
			BeforeEach(func() {
				scope = "/aaa/bbb/ccc/"
			})

			It("the returned scope should be the same as the original", func() {
				Expect(newScope).To(Equal(scope))
			})
		})
		Context("scope is not a sub resource", func() {
			BeforeEach(func() {
				scope = "/aaa/"
			})

			It("the returned scope should be the same as the original", func() {
				Expect(newScope).To(Equal(scope))
			})
		})
	})
	Context("json string is valid", func() {
		BeforeEach(func() {
			pathFmtJson = `{"api": "/%s/aaa"}`
			subPathFmtJson = `{"api": "/%s/aaa/%s/"}`
		})
		When("the specified scene is not exist", func() {
			BeforeEach(func() {
				scene = "not-exist"
			})
			Context("scope is a sub resource", func() {
				BeforeEach(func() {
					scope = "/aaa/bbb/ccc/"
				})

				It("the returned scope should be the same as the original", func() {
					Expect(newScope).To(Equal(scope))
				})
			})
			Context("scope is not a sub resource", func() {
				BeforeEach(func() {
					scope = "/aaa/"
				})

				It("the returned scope should be the same as the original", func() {
					Expect(newScope).To(Equal(scope))
				})
			})
		})
		When("the specified scene exists", func() {
			BeforeEach(func() {
				scene = "api"
			})
			Context("scope is a sub resource", func() {
				BeforeEach(func() {
					scope = "/111/222/333/"
				})

				It("the returned scope is as expected", func() {
					Expect(newScope).To(Equal("/111/aaa/222/333/"))
				})
			})
			Context("scope is not a sub resource", func() {
				BeforeEach(func() {
					scope = "/111/"
				})

				It("the returned scope is as expected", func() {
					Expect(newScope).To(Equal("/111/aaa/"))
				})
			})
		})
	})
})

var _ = Describe("TestResourcePathFormat_FormatPathAllScene", func() {
	var (
		pathFmtJson    string
		subPathFmtJson string
		scope          string
		newScopes      []string

		r *ResourcePathFormat
	)

	BeforeEach(func() {
		pathFmtJson = ""
		subPathFmtJson = ""
		scope = "/"
		newScopes = []string{}
		r = &ResourcePathFormat{}
	})

	JustBeforeEach(func() {
		r = NewResourcePathFormat(pathFmtJson, subPathFmtJson)
		newScopes = r.FormatPathAllScene(scope)
	})
	Context("json string is invalid", func() {
		BeforeEach(func() {
			pathFmtJson = "invalid json"
			subPathFmtJson = "invalid json"
		})
		Context("scope is a sub resource", func() {
			It("the returned scope should be the same as the original", func() {
				Expect(newScopes).To(Equal([]string{scope}))
			})
		})
		Context("scope is not a sub resource", func() {
			It("the returned scope should be the same as the original", func() {
				Expect(newScopes).To(Equal([]string{scope}))
			})
		})
	})
	Context("json string is valid", func() {
		BeforeEach(func() {
			pathFmtJson = `{"api": "/%s/aaa/", "web": "/%s/web/"}`
			subPathFmtJson = `{"api": "/%s/aaa/%s/", "web": "/%s/web/%s/"}`
		})
		Context("scope is a sub resource", func() {
			BeforeEach(func() {
				scope = "/111/222/333/"
			})

			It("the returned scope is as expected", func() {
				Expect(newScopes).To(ContainElements(
					"/111/aaa/222/333/",
					"/111/web/222/333/",
				))
			})
		})
		Context("scope is not a sub resource", func() {
			BeforeEach(func() {
				scope = "/111/"
			})

			It("the returned scope is as expected", func() {
				Expect(newScopes).To(ContainElements(
					"/111/aaa/",
					"/111/web/",
				))
			})
		})
	})
})
