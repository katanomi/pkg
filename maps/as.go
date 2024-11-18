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
	"bytes"
	"fmt"
	"strconv"

	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// AsBool will parse value in data["key"] to bool
// if the key not exist, it will return nil
// if the value is not a valid bool string, it will return error
func AsBool(data map[string]string, key string) (*bool, error) {
	value, ok := data[key]
	if !ok {
		return nil, nil
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		falseB := false
		return &falseB, fmt.Errorf("parse key '%s' error: %s, value: %s", key, err.Error(), value)
	}
	return &v, nil
}

// AsObject will parse value in data["key"] to object
// it expects the value should be yaml content
func AsObject(data map[string]string, key string, object interface{}) error {
	value, ok := data[key]
	if !ok {
		return nil
	}

	err := utilyaml.NewYAMLToJSONDecoder(bytes.NewBufferString(value)).Decode(object)
	if err != nil {
		return fmt.Errorf("parse key '%s' error: %s, value: %s", key, err.Error(), value)
	}
	return nil
}
