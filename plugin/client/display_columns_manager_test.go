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

package client

import (
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.DisplayColumnsManager", func() {
	var (
		displayColumnsManager *DisplayColumnsManager
	)
	BeforeEach(func() {
		displayColumnsManager = &DisplayColumnsManager{}
	})
	Context("set display columns", func() {
		BeforeEach(func() {
			displayColumnsManager.SetDisplayColumns("key", v1alpha1.DisplayColumn{Name: "value"})
		})

		It("should get display columns", func() {
			Expect(displayColumnsManager.GetDisplayColumns()).To(Equal(map[string]v1alpha1.DisplayColumns{"key": {{Name: "value"}}}))
		})
	})
})
