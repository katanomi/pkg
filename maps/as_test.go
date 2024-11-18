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

var data = map[string]string{
	"stringempty": "",
	"string":      "a1",
	"bool":        "false",
	"array": `
- a1
- b1
`,
	"object": `
key1: value1
key2: value2
`,
	"object2": `
key1: value1
key2: |
- v2.1
- v2.2
`,
}

type obj struct {
	Key1 string `json:"key1,omitempty"`
	Key2 string `json:"key2,omitempty"`
}

func TestAsBool(t *testing.T) {
	t.Run("key is not existed, it should return nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		res, err := AsBool(data, "what")
		g.Expect(err).Should(gomega.BeNil())
		g.Expect(res).Should(gomega.BeNil())
	})

	t.Run("key is existed, value is valid, it should parse correctly", func(t *testing.T) {
		g := gomega.NewWithT(t)

		res, err := AsBool(data, "bool")
		g.Expect(err).Should(gomega.BeNil())
		g.Expect(*res).Should(gomega.BeFalse())
	})

	t.Run("key is existed, value is invalid, it should return error", func(t *testing.T) {
		g := gomega.NewWithT(t)

		_, err := AsBool(data, "stringempty")
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestAsObject(t *testing.T) {

	t.Run("key is not existed, it should do nothing with object", func(t *testing.T) {
		g := gomega.NewWithT(t)

		o := &obj{}
		err := AsObject(data, "what", o)
		g.Expect(err).Should(gomega.BeNil())
		g.Expect(o.Key1).Should(gomega.BeEmpty())
		g.Expect(o.Key2).Should(gomega.BeEmpty())
	})

	t.Run("key is existed, but the value is invalid, it should return error", func(t *testing.T) {
		g := gomega.NewWithT(t)

		o := &obj{}
		err := AsObject(data, "string", o)
		g.Expect(err).ShouldNot(gomega.BeNil())
	})

	t.Run("key is existed, and the value is valid object, it should return correct object", func(t *testing.T) {
		g := gomega.NewWithT(t)

		o := &obj{}
		err := AsObject(data, "object", o)
		g.Expect(err).Should(gomega.BeNil())
		g.Expect(o.Key1).Should(gomega.BeEquivalentTo("value1"))
		g.Expect(o.Key2).Should(gomega.BeEquivalentTo("value2"))
	})
	t.Run("key is existed, and the value is valid object, but the object type is not matched, it should return error", func(t *testing.T) {
		g := gomega.NewWithT(t)

		o := &obj{}
		err := AsObject(data, "object2", o)
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}
