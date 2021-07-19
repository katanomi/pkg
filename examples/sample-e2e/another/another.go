// +build e2e

/*
Copyright 2021 The Katanomi Authors.

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

package another

import (
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = P1Case("another-test").Cluster().WithFunc(func(ctx TestContext) {
	// test case
	BeforeEach(func() {
		ctx.Debugw("before each in another pkg")
	})
	AfterEach(func() {
		ctx.Debugw("after each in another pkg")
	})
	Context("With a cluster scoped test case", func() {
		JustBeforeEach(func() {
			ctx.Infow("just before each in another pkg")
		})
		JustAfterEach(func() {
			ctx.Infow("just after each in another pkg")
		})
		It("it", func() {
			Expect(ctx.Config).ToNot(BeNil())
		})
	})
}).Do()
