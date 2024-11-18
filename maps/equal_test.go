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

package maps

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestSameStringMap(t *testing.T) {

	gomega.RegisterTestingT(t)
	t.Run("same string map", func(t *testing.T) {
		first := map[string]string{
			"a": "b",
			"b": "a",
		}
		second := map[string]string{
			"b": "a",
			"a": "b",
		}
		gomega.Expect(IsSameStringMap(first, second)).To(gomega.BeTrue())

	})

	t.Run("different string map", func(t *testing.T) {
		first := map[string]string{
			"a": "d",
			"b": "a",
		}
		second := map[string]string{
			"b": "c",
			"a": "b",
		}
		gomega.Expect(IsSameStringMap(first, second)).To(gomega.BeFalse())

	})

	t.Run("different length string map", func(t *testing.T) {
		first := map[string]string{
			"a": "d",
			"b": "a",
			"c": "e",
		}
		second := map[string]string{
			"b": "c",
			"a": "b",
		}
		gomega.Expect(IsSameStringMap(first, second)).To(gomega.BeFalse())

	})
}
