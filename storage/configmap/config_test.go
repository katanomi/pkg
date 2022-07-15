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

package configmap

import (
	"sync/atomic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("testing for config", func() {
	When("configmap is nil", func() {
		It("", func() {
			testCM := testConfigMap("cc", nil)
			testCount := int64(0)
			c := NewConfigConstructor(nil, func(configMap *corev1.ConfigMap) {
				atomic.AddInt64(&testCount, 1)
				Expect(configMap.GetName()).To(Equal("cc"))
			})
			Expect(c.CmName()).To(Equal(""))
			Expect(c.Default()).To(BeNil())
			c.Handle(testCM)
			Expect(testCount).To(Equal(int64(1)))
		})
	})
})
