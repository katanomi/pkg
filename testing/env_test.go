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

package testing

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDefaultEnv", func() {
	key := "test-env"
	defaultValue := "default-value"
	When("provide variables through environment variables", func() {
		envValue := "env-value"
		BeforeEach(func() {
			err := os.Setenv(key, envValue)
			Expect(err).To(Succeed())

			DeferCleanup(func() {
				err := os.Unsetenv(key)
				Expect(err).To(Succeed())
			})
		})

		It("the value should equal to the envValue", func() {
			value := GetDefaultEnv(key, defaultValue)
			Expect(value).To(Equal(envValue))
		})
	})

	When("not provide environment variable", func() {
		It("should get the default value", func() {
			value := GetDefaultEnv(key, defaultValue)
			Expect(value).To(Equal(defaultValue))
		})
	})
})
