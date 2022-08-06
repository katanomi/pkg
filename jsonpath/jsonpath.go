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

// Package jsonpath provide readablity and writeablity of golang struct by jsonpath
package jsonpath

import (
	"fmt"
	"reflect"

	"github.com/yuzp1996/client-go/util/jsonpath"
)

// Read will read value of data by jsonpath
func Read(data interface{}, path string) ([][]interface{}, error) {
	jpath := jsonpath.New("jsonpath")
	jpath.AllowMissingKeys(true)
	err := jpath.Parse(path)
	if err != nil {
		return nil, err
	}
	fullResults, err := jpath.FindResults(data)
	if err != nil {
		return nil, err
	}

	var readRes = [][]interface{}{}
	for i := range fullResults {
		items := fullResults[i]
		readResItems := []interface{}{}
		for j := range items {
			resultValue := items[j]
			readResItems = append(readResItems, resultValue.Interface())
		}
		readRes = append(readRes, readResItems)
	}

	return readRes, nil
}

// Write will set value of data by jsonpath, the data must be pointer
func Write(data interface{}, path string, value interface{}) error {
	t := reflect.ValueOf(data)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("data should be Pointer, but: %T", data)
	}

	jpath := jsonpath.New("jsonpath")
	jpath.AllowMissingKeys(true)
	err := jpath.Parse(path)
	if err != nil {
		return err
	}
	fullResults, err := jpath.FindResults(data)
	if err != nil {
		return err
	}

	for i := range fullResults {
		for j := range fullResults[i] {
			f := fullResults[i][j]
			f.Set(reflect.ValueOf(value))
		}
	}

	return nil
}
