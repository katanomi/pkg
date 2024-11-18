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

package testing

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
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

func TestLoadMultiYaml_success(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(LoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.success.yaml", &cms)).Should(BeNil())
	g.Expect(cms).To(HaveLen(2))
	g.Expect(cms[0].Name).To(Equal("abc-1"))
	g.Expect(cms[1].Name).To(Equal("abc-2"))
}

func TestLoadMultiJson_success(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []*unstructured.Unstructured{}
	g.Expect(LoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.success.json", &cms)).Should(BeNil())
	g.Expect(cms).To(HaveLen(2))
	g.Expect(cms[0].GetName()).To(Equal("abc-1"))
	g.Expect(cms[1].GetName()).To(Equal("abc-2"))
}

func TestLoadMultiJson_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(LoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.fail.json", &cms)).ShouldNot(BeNil())
	g.Expect(cms).To(BeEmpty())
}

func TestLoadMultiYaml_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(LoadMultiYamlOrJson("./testdata/not-exist.yaml", &cms)).ShouldNot(BeNil())
}

func TestMustLoadMultiYaml_success(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(func() {
		MustLoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.success.yaml", &cms)
	}).ShouldNot(Panic())
}

func TestMustLoadMultiJson_success(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(func() {
		MustLoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.success.json", &cms)
	}).ShouldNot(Panic())
}

func TestMustLoadMultiJson_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(func() {
		MustLoadMultiYamlOrJson("./testdata/loadMultiYaml.configmap.fail.json", &cms)
	}).Should(Panic())
}

func TestMustLoadMultiYaml_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	cms := []corev1.ConfigMap{}
	g.Expect(func() {
		MustLoadMultiYamlOrJson("./testdata/not-exist.yaml", &cms)
	}).Should(Panic())
}

func TestMustLoadFileString_success(t *testing.T) {
	g := NewGomegaWithT(t)
	var strs string
	g.Expect(func() {
		MustLoadFileString("./testdata/loadMultiYaml.configmap.success.json", &strs)
	}).ShouldNot(Panic())
	g.Expect(strs).NotTo(BeEmpty())
}

func TestMustLoadFileString_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	var strs string
	g.Expect(func() {
		MustLoadFileString("./testdata/not-exist.yaml", &strs)
	}).Should(Panic())
}

func TestMustLoadFileBytes_success(t *testing.T) {
	g := NewGomegaWithT(t)
	var content []byte
	g.Expect(func() {
		content = MustLoadFileBytes("./testdata/loadMultiYaml.configmap.success.json")
	}).ShouldNot(Panic())
	g.Expect(content).NotTo(BeEmpty())
}

func TestMustLoadFileBytes_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(func() {
		MustLoadFileBytes("./testdata/not-exist.yaml")
	}).Should(Panic())
}
