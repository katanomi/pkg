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

package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.ApplyControllerBuilderOptions", func() {
	builder := &builder.Builder{}
	builderFuncs := make([]ControllerBuilderOption, 0)

	JustBeforeEach(func() {
		ApplyControllerBuilderOptions(builder, builderFuncs...)
	})

	BeforeEach(func() {
		builderFuncs = append(builderFuncs, WithBuilderOptions(controller.Options{
			RateLimiter: DefaultTypedRateLimiter[reconcile.Request](),
		}))
	})

	It("builds options", func() {
		Expect(builder).NotTo(BeNil())
	})
})
