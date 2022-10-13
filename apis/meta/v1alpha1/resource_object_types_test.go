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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/katanomi/pkg/testing"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.ResourceURL.Validate", func() {

	var (
		resourceURL *ResourceURL
		errs        field.ErrorList
		path        *field.Path
	)

	BeforeEach(func() {
		resourceURL = &ResourceURL{}
		errs = nil
		path = field.NewPath("spec")
	})

	JustBeforeEach(func() {
		errs = resourceURL.Validate(path)
	})

	Context("empty resource url", func() {
		It("should return an error", func() {
			Expect(errs).ToNot(BeNil())
			Expect(errs).To(HaveLen(1))
		})
	})

	Context("valid resource url", func() {
		BeforeEach(func() {
			MustLoadYaml("testdata/resource_url.valid.yaml", resourceURL)
		})
		It("should not return an error", func() {
			Expect(errs).To(HaveLen(0))
		})
	})

})
