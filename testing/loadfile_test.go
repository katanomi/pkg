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

package testing

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMustLoadJSON_success(t *testing.T) {
	g := NewGomegaWithT(t)
	m := make(map[string]string)
	g.Expect(func() {
		MustLoadJSON("./testdata/valid_json.json", &m)
	}).ShouldNot(Panic())
	g.Expect(m).NotTo(BeNil())
	g.Expect(m["name"]).To(Equal("tom"))
}

func TestMustLoadJSON_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	m := make(map[string]string)
	g.Expect(func() {
		MustLoadJSON("./testdata/invalid_json.json", &m)
	}).Should(Panic())
	g.Expect(m).To(HaveLen(0))
}
