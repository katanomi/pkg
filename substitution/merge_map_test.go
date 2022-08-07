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

package substitution

import (

	// ktesting "github.com/katanomi/pkg/testing"
	"bytes"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// corev1 "k8s.io/api/core/v1"

	"knative.dev/pkg/logging"
)

func DeepCopyMap(dst, src *map[string]string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

var _ = Describe("Test.MergeMap", func() {
	var (
		left, right, oleft, oright, merge, actual map[string]string
		//
		// log level. It can be debug, info, warn, error, dpanic, panic, fatal.
		log, _ = logging.NewLogger("", "debug")
	)

	BeforeEach(func() {
		left = map[string]string{
			"l1": "v1",
			"l2": "v2",
			"l3": "v3",
		}
		right = map[string]string{
			"r1": "v1",
			"r2": "v2",
			"r3": "v3",
		}
		merge = map[string]string{
			"l1": "v1",
			"l2": "v2",
			"l3": "v3",
			"r1": "v1",
			"r2": "v2",
			"r3": "v3",
		}
		DeepCopyMap(&oleft, &left)
		DeepCopyMap(&oright, &right)
	})

	JustBeforeEach(func() {
		log.Infow("MergeMap", "left", left, "right", right)
		actual = MergeMap(left, right)
	})

	When("left is empty", func() {
		BeforeEach(func() {
			left = nil
		})

		It("should merge success", func() {
			Expect(actual).To(Equal(oright))
		})
	})

	When("right is empty", func() {
		BeforeEach(func() {
			right = nil
		})

		It("should merge success", func() {
			Expect(actual).To(Equal(oleft))
		})
	})

	When("there is no common key", func() {
		It("should merge success", func() {
			Expect(actual).To(Equal(merge))
		})
	})

	When("there are common keys", func() {
		BeforeEach(func() {
			right = left
		})
		It("should merge success", func() {
			Expect(actual).To(Equal(oleft))
		})
	})
})
