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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"yunion.io/x/pkg/errors"
)

func LoadKubeResourcesAsUnstructured(file string) (objs []unstructured.Unstructured, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}
	objs = []unstructured.Unstructured{}
	parts := strings.Split(string(data), "---")
	for _, y := range parts {
		if len(strings.TrimSpace(y)) == 0 {
			continue
		}
		obj := &unstructured.Unstructured{}
		err = utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(y)), len([]byte(y))).Decode(obj)
		if err != nil {
			return
		}
		if obj != nil {
			objs = append(objs, *obj)
		}
	}
	return
}

// LoadKubeResources loading kubernetes resources
func LoadKubeResources(file string, clt client.Client) (err error) {
	objs, err := LoadKubeResourcesAsUnstructured(file)
	if err != nil {
		return
	}
	for _, obj := range objs {
		if err = clt.Create(context.Background(), &obj); err != nil {
			return
		}
	}
	return
}

// UnstructedToTyped converts an unstructured object into a object
// Warning: This SHOULD never be used in production code, only in test code
func UnstructedToTyped(from unstructured.Unstructured, to interface{}) error {
	data, err := json.Marshal(from.Object)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, to)
}

// DeleteResources delete resources contained in the file
func DeleteResources(file string, clt client.Client) (err error) {
	objs, err := LoadKubeResourcesAsUnstructured(file)
	if err != nil {
		return err
	}
	errs := []error{}
	for _, obj := range objs {
		if err = clt.Delete(context.Background(), &obj); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.NewAggregate(errs)
}
