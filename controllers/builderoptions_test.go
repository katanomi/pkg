/*
Copyright 2023 The AlaudaDevops Authors.

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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Test.BuilderOptions", func() {
	var (
		opts         controller.Options
		buildOptFuns []BuilderOptionFunc
	)

	BeforeEach(func() {
		opts = DefaultOptions()
	})
	JustBeforeEach(func() {
		opts = BuilderOptions(buildOptFuns...)
	})

	Context("BuilderOptions empty", func() {
		It("should return default options", func() {
			Expect(opts).To(Equal(DefaultOptions()))
		})
	})

	Context("builderoptions with custom rate limiter", func() {
		customRateLimiter := workqueue.NewTypedMaxOfRateLimiter[reconcile.Request]()

		BeforeEach(func() {
			buildOptFuns = append(buildOptFuns, RateLimiter(customRateLimiter))
		})
		It("should return options with custom rateLimiter", func() {
			Expect(opts.RateLimiter).To(Equal(customRateLimiter))
		})
	})

	Context("builderoptions with custom maxConcurrentReconciles", func() {
		BeforeEach(func() {
			buildOptFuns = append(buildOptFuns, MaxConCurrentReconciles(100))
		})
		It("should return options with custom rateLimiter", func() {
			Expect(opts.MaxConcurrentReconciles).To(Equal(100))
		})
	})

})
