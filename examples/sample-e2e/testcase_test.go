//go:build e2e
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

package e2e

import (
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = TestCase(Options{Name: "Testinge2e", Priority: P0, Scope: NamespaceScoped}).WithFunc(func(ctx TestContext) {
	BeforeEach(func() {
		ctx.Debugw("some debug message")
		// fmt.Println("TestCase BeforeEach", ctx.Config)
	})
	It("should succeed", func() {
		Expect(ctx.Config).ToNot(BeNil())
	})
}).Do()

var _ = P0Case("aaanother-test").Cluster().WithFunc(func(ctx TestContext) {
	// test case
	It("should succeed", func() {
		Expect(ctx.Config).ToNot(BeNil())
	})
}).Do()
