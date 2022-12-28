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

package encoding

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var DefaultJsonPath = JsonPath{}

// Encode using DefaultJsonPath encode obj to json path
func Encode(obj interface{}) map[string]string {
	return DefaultJsonPath.Encode(obj)
}

// Encode using DefaultJsonPath decode json path to obj
func Decode(obj interface{}, data map[string]string) error {
	return DefaultJsonPath.Decode(obj, data)
}

// NewJsonPath construct json path
func NewJsonPath() JsonPath {
	return JsonPath{}
}

// JsonPath object for json-path encoding and decoding
type JsonPath struct {
	// PathFormat is the format of the path
	PathFormat func(string) string
}

func (p JsonPath) formatPath(s string) string {
	if p.PathFormat != nil {
		return p.PathFormat(s)
	}

	return s
}

// Decode decode json path to obj
func (p JsonPath) Decode(obj interface{}, data map[string]string) error {
	root := NewRootNode()
	for k, v := range data {
		assignNodeValue(strings.Split(k, "."), v, root)
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(root.Json()), &m); err != nil {
		return err
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           obj,
		TagName:          "path",
	}

	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(m)
}

// Encode encode obj to json path
func (p JsonPath) Encode(obj interface{}) map[string]string {
	value := reflect.ValueOf(obj)
	m := make(map[string]string)
	return p.jsonPathEncode(value, "", m)
}

func (p JsonPath) getStructFieldName(field reflect.StructField) string {
	fname := p.formatPath(field.Name)
	if tagName := field.Tag.Get("json"); tagName != "" {
		fname = strings.SplitN(tagName, ",", 2)[0]
	}
	if tagName := field.Tag.Get("path"); tagName != "" {
		fname = strings.SplitN(tagName, ",", 2)[0]
	}

	return fname
}

func (p JsonPath) jsonPathEncode(t reflect.Value, name string, dst map[string]string) map[string]string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		return p.jsonPathEncode(t.Elem(), name, dst)
	case reflect.Map:
		for _, key := range t.MapKeys() {
			fname := p.formatPath(fmt.Sprintf("%v", key.Interface()))
			fname = strings.Trim(name+"."+fname, ".")
			dst = p.jsonPathEncode(t.MapIndex(key), fname, dst)
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fname := p.getStructFieldName(t.Type().Field(i))
			fname = strings.Trim(name+"."+fname, ".")
			dst = p.jsonPathEncode(t.Field(i), fname, dst)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < t.Len(); i++ {
			dst = p.jsonPathEncode(t.Index(i), name+"["+strconv.Itoa(i)+"]", dst)
		}

	default:
		dst[name] = fmt.Sprintf("%v", t.Interface())
	}
	return dst
}
