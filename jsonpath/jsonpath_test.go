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

package jsonpath

import (
	"reflect"
	"testing"
)

type book struct {
	Name     string        `json:"name"`
	Price    int           `json:"price"`
	Sections []bookSection `json:"sections"`
}

type bookSection struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

var data book = book{
	Name:  "katanomi",
	Price: 100,
	Sections: []bookSection{
		{
			Name:  "pkg",
			Title: "pkg description",
		},
		{
			Name:  "core",
			Title: "core description",
		},
		{
			Name:  "builds",
			Title: "builds description",
		},
	},
}

func TestRead(t *testing.T) {
	type caseData struct {
		Description   string
		Path          string
		Expected      interface{}
		ExpectedError bool
	}

	var items = []caseData{
		{"error field path", "{.name[]}", "", true},
		{"field path of struct, type is string", "{.name}", [][]interface{}{{"katanomi"}}, false},
		{"field path of struct, type is int", "{.price}", [][]interface{}{{int(100)}}, false},
		{"index filter in array object", "{.sections[1].name}", [][]interface{}{{"core"}}, false},
		{"attribute filter in array object", "{.sections[?(@.name=='pkg')].title}", [][]interface{}{{"pkg description"}}, false},
		{"attribute filter in array object and return array", "{.sections[?(@.name!='pkg')].title}", [][]interface{}{{"core description", "builds description"}}, false},
	}

	for _, item := range items {
		t.Run(item.Description, func(t *testing.T) {
			vals, err := Read(data, item.Path)
			if item.ExpectedError {
				if err == nil {
					t.Errorf("err should not nil")
					return
				}
				return
			}
			if !reflect.DeepEqual(vals, item.Expected) {
				t.Errorf("expected: %#v, but: %#v", item.Expected, vals)
			}
		})
	}
}

func TestWrite(t *testing.T) {

	type caseData struct {
		Description string
		Path        string
		SetValue    interface{}
		Assert      func(dataW *book) bool
	}

	var items = []caseData{
		{
			"field path of struct, type is string",
			"{.name}",
			"katanomi-1",
			func(dataW *book) bool { return dataW.Name == "katanomi-1" }},
		{
			"field path of struct, type is int",
			"{.price}",
			int(1000),
			func(dataW *book) bool {
				return dataW.Price == 1000
			},
		},
		{
			"index filter in array object",
			"{.sections[1].name}",
			"core-1",
			func(dataW *book) bool {
				return dataW.Sections[1].Name == "core-1"
			},
		},
		{
			"attribute filter in array object",
			"{.sections[?(@.name=='pkg')].title}",
			"pkg description 1",
			func(dataW *book) bool {
				return dataW.Sections[0].Title == "pkg description 1"
			},
		},
		{
			"attribute filter in array object and return array",
			"{.sections[?(@.name!='pkg')].title}",
			"not pkg description",
			func(dataW *book) bool {
				return dataW.Sections[1].Title == "not pkg description" && dataW.Sections[2].Title == "not pkg description"
			},
		},
	}

	for _, item := range items {
		t.Run(item.Description, func(t *testing.T) {
			var dataW = data
			Write(&dataW, item.Path, item.SetValue)
			if !item.Assert(&dataW) {
				t.Errorf("could not write correctly, data: %#v", dataW)
			}
		})
	}

	err := Write(data, "{.what}", "xx")
	if err == nil {
		t.Errorf("write data must be pointer, err should not nil")
	}
}
