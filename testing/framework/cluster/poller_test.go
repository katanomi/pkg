/*
Copyright 2022 The AlaudaDevops Authors.

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

package cluster

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Poller Settings", func() {
	var poller *Poller

	BeforeEach(func() {
		poller = &Poller{}
	})

	Context("provide `interval` and `timeout` parameter", func() {
		When("the value of `interval` is greater than `timeout`", func() {
			BeforeEach(func() {
				poller.Interval = time.Second * 5
				poller.Timeout = time.Second * 3
			})

			It("the value of `timeout` should be 5", func() {
				interval, timeout := poller.Settings()
				Expect(interval).To(Equal(time.Second * 5))
				Expect(timeout).To(Equal(time.Second * 5))
			})
		})
	})

	Context("not provide `interval` and `timeout` parameter", func() {
		It("should be the default value", func() {
			interval, timeout := poller.Settings()
			Expect(interval).To(Equal(200 * time.Millisecond))
			Expect(timeout).To(Equal(5 * time.Second))
		})
	})
})
