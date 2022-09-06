/*
Copyright 2021 The Katanomi Authors.

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
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

// MustLoadFileString loads a file as string
// will panic if if failes
// ONLY FOR TEST USAGE
func MustLoadFileString(file string, content *string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	*content = string(data)
}

// LoadJSON loads json
func LoadJSON(file string, obj interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}

// MustLoadJSON loads json or panics if the parse fails.
func MustLoadJSON(file string, obj interface{}) {
	err := LoadJSON(file, obj)
	if err != nil {
		panic(fmt.Sprintf("load json file failed, file path: %s, err: %s", file, err))
	}
}

// LoadYAML loads yaml
func LoadYAML(file string, obj interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}
	err = yaml.Unmarshal(data, obj)
	return
}

// MustLoadYaml loads yaml or panics if the parse fails.
func MustLoadYaml(file string, obj interface{}) {
	err := LoadYAML(file, obj)
	if err != nil {
		panic(fmt.Sprintf("load yaml file failed, file path: %s, err: %s", file, err))
	}
}

// LoadObjectOrDie loads object from yaml and returns
func LoadObjectOrDie(g *WithT, file string, obj metav1.Object, patches ...func(metav1.Object)) metav1.Object {
	g.Expect(LoadYAML(file, obj)).To(Succeed(), "could not load file into metav1.Object")
	for _, p := range patches {
		p(obj)
	}
	return obj
}

// LoadObjectReferenceOrDie loads object reference from yaml and returns
func LoadObjectReferenceOrDie(g *WithT, file string, obj *corev1.ObjectReference, patches ...func(*corev1.ObjectReference)) *corev1.ObjectReference {
	g.Expect(LoadYAML(file, obj)).To(Succeed(), "could not load file into corev1.ObjectReference")
	for _, p := range patches {
		p(obj)
	}
	return obj
}

// MustLoadReturnObjectFromYAML loads and object from yaml file and returns as interface{}
// if any loading errors happen will panic
// TO BE USED IN TESTS, DO NOT USE IN PRODUCTION CODE
func MustLoadReturnObjectFromYAML(file string, obj interface{}) interface{} {
	MustLoadYaml(file, obj)
	return obj
}
