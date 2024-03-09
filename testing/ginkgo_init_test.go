/*
Copyright 2024 The Katanomi Authors.

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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"go.uber.org/zap"
)

var _ = Describe("Test.InitGinkgoWithLogger", func() {
	var logger *zap.SugaredLogger
	BeforeEach(func() {
		logger = InitGinkgoWithLogger()
	})
	It("should disable string length limit in Gomega", func() {
		Expect(format.MaxLength).To(Equal(0))
	})
	It("should return a logger", func() {
		Expect(logger).NotTo(BeNil())
	})
})
