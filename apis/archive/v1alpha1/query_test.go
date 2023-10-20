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

package v1alpha1

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
)

func checkCondition(file string, t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	cond := Condition{}
	ktesting.MustLoadYaml(file, &cond)
	data, _ := json.Marshal(cond)
	got := map[string]interface{}{}
	_ = json.Unmarshal(data, &got)

	want := map[string]interface{}{}
	ktesting.MustLoadYaml(file, &want)

	g.Expect(cmp.Diff(got, want)).To(gomega.BeEmpty())
}

func TestCondition(t *testing.T) {
	checkCondition("testdata/condition.or.golden.yaml", t)
	checkCondition("testdata/condition.eq.column.golden.yaml", t)
	checkCondition("testdata/condition.and.golden.yaml", t)
	checkCondition("testdata/condition.in.string.golden.yaml", t)
}
