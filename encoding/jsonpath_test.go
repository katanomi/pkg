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

package encoding

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	testing2 "github.com/AlaudaDevops/pkg/testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testStruct struct {
	Int     int      `json:"int"`
	String  string   `json:"string"`
	Float64 float64  `json:"float64"`
	Slice   []string `json:"slice_path" path:"slice_path"`
}

func Test_JsonPathEncode_struct(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		obj   interface{}
		paths map[string]string
	}{
		{
			obj: testStruct{
				Int:     1,
				String:  "s",
				Float64: 1.1,
				Slice:   []string{"1.1", "1.2"},
			},
			paths: map[string]string{
				"int":           "1",
				"string":        "s",
				"float64":       "1.1",
				"slice_path[0]": "1.1",
				"slice_path[1]": "1.2",
			},
		},
		{
			obj: []testStruct{
				{
					Int:     1,
					String:  "s",
					Float64: 1.1,
					Slice:   []string{"1.1", "1.2"},
				},
				{
					Int:     2,
					String:  "s2",
					Float64: 2.1,
					Slice:   []string{"2.1", "2.2"},
				},
			},
			paths: map[string]string{
				"[0].int":           "1",
				"[0].string":        "s",
				"[0].float64":       "1.1",
				"[0].slice_path[0]": "1.1",
				"[0].slice_path[1]": "1.2",
				"[1].int":           "2",
				"[1].string":        "s2",
				"[1].float64":       "2.1",
				"[1].slice_path[0]": "2.1",
				"[1].slice_path[1]": "2.2",
			},
		},
	}

	jsonpath := NewJsonPath()
	for _, tt := range tests {
		g.Expect(jsonpath.Encode(tt.obj)).To(Equal(tt.paths))
	}
}

func Test_JsonPathEncode_map(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		obj   map[string]interface{}
		paths map[string]string
	}{
		{
			obj: map[string]interface{}{
				"str_key":   "value1",
				"int_key":   1,
				"float_key": 1.234,
			},
			paths: map[string]string{
				"str_key":   "value1",
				"int_key":   "1",
				"float_key": "1.234",
			},
		},
		{
			obj: map[string]interface{}{
				"str_slice_key":   []string{"value1", "value2"},
				"int_slice_key":   []int{1, 2, 3},
				"float_slice_key": []float64{1.1, 2.2, 3.3},
			},
			paths: map[string]string{
				"str_slice_key[0]":   "value1",
				"str_slice_key[1]":   "value2",
				"int_slice_key[0]":   "1",
				"int_slice_key[1]":   "2",
				"int_slice_key[2]":   "3",
				"float_slice_key[0]": "1.1",
				"float_slice_key[1]": "2.2",
				"float_slice_key[2]": "3.3",
			},
		},
		{
			obj: map[string]interface{}{
				"struct_key": testStruct{
					Int:     1,
					String:  "test-string",
					Float64: 1.234,
					Slice:   []string{"value1", "value2"},
				},
			},
			paths: map[string]string{
				"struct_key.int":           "1",
				"struct_key.string":        "test-string",
				"struct_key.float64":       "1.234",
				"struct_key.slice_path[0]": "value1",
				"struct_key.slice_path[1]": "value2",
			},
		},
		{
			obj: map[string]interface{}{
				"struct_slice_key": []testStruct{
					{
						Int:     1,
						String:  "test-string",
						Float64: 1.234,
						Slice:   []string{"value1", "value2"},
					},
					{
						Int:     2,
						String:  "test-string2",
						Float64: 2.234,
						Slice:   []string{"value21", "value22"},
					},
				},
			},
			paths: map[string]string{
				"struct_slice_key[0].int":           "1",
				"struct_slice_key[0].string":        "test-string",
				"struct_slice_key[0].float64":       "1.234",
				"struct_slice_key[0].slice_path[0]": "value1",
				"struct_slice_key[0].slice_path[1]": "value2",
				"struct_slice_key[1].int":           "2",
				"struct_slice_key[1].string":        "test-string2",
				"struct_slice_key[1].float64":       "2.234",
				"struct_slice_key[1].slice_path[0]": "value21",
				"struct_slice_key[1].slice_path[1]": "value22",
			},
		},
	}

	jsonpath := NewJsonPath()
	for _, tt := range tests {
		g.Expect(jsonpath.Encode(tt.obj)).To(Equal(tt.paths))
	}
}

var _ = Describe("TestJsonPath_Decode", func() {
	var m map[string]string

	Context("test decode struct", func() {
		BeforeEach(func() {
			m = make(map[string]string)
			testing2.MustLoadJSON("./testdata/struct_jsonpath.json", &m)
		})

		It("should decode to struct successfully", func() {
			type testObj struct {
				StructSliceKey []testStruct `json:"struct_slice_key" path:"struct_slice_key"`
			}
			obj := testObj{}
			err := NewJsonPath().Decode(&obj, m)
			Expect(err).Should(Succeed())

			goldenData := testObj{}
			testing2.MustLoadJSON("./testdata/struct_jsonpath_golden.json", &goldenData)
			diff := cmp.Diff(obj, goldenData)
			Expect(diff).To(BeEmpty())
		})
	})

	Context("test decode struct with inline struct", func() {
		BeforeEach(func() {
			m = make(map[string]string)
			testing2.MustLoadJSON("./testdata/struct_inline_jsonpath.json", &m)
		})

		It("should decode to struct successfully", func() {
			type Person struct {
				Name string
			}
			type testObj struct {
				Person         `path:",squash"`
				StructSliceKey []testStruct `json:"struct_slice_key" path:"struct_slice_key"`
			}
			obj := testObj{}
			err := NewJsonPath().Decode(&obj, m)
			Expect(err).Should(Succeed())

			goldenData := testObj{}
			testing2.MustLoadJSON("./testdata/struct_inline_jsonpath_golden.json", &goldenData)
			diff := cmp.Diff(obj, goldenData)
			Expect(diff).To(BeEmpty())
		})
	})
})
