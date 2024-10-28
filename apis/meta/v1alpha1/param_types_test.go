/*
Copyright 2024 The Katanomi Authors.

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

package v1alpha1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	ktesting "github.com/katanomi/pkg/testing"

	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/util/sets"
	"knative.dev/pkg/apis"
)

func TestParamSpec_SetDefaults(t *testing.T) {
	tests := []struct {
		name            string
		before          *v1alpha1.ParamSpec
		defaultsApplied *v1alpha1.ParamSpec
	}{{
		name: "inferred string type",
		before: &v1alpha1.ParamSpec{
			Name: "parametername",
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name: "parametername",
			Type: v1alpha1.ParamTypeString,
		},
	}, {
		name: "inferred type from default value - array",
		before: &v1alpha1.ParamSpec{
			Name: "parametername",
			Default: &v1alpha1.ParamValue{
				ArrayVal: []string{"array"},
			},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name: "parametername",
			Type: v1alpha1.ParamTypeArray,
			Default: &v1alpha1.ParamValue{
				ArrayVal: []string{"array"},
			},
		},
	}, {
		name: "inferred type from default value - string",
		before: &v1alpha1.ParamSpec{
			Name: "parametername",
			Default: &v1alpha1.ParamValue{
				StringVal: "an",
			},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name: "parametername",
			Type: v1alpha1.ParamTypeString,
			Default: &v1alpha1.ParamValue{
				StringVal: "an",
			},
		},
	}, {
		name: "inferred type from default value - object",
		before: &v1alpha1.ParamSpec{
			Name: "parametername",
			Default: &v1alpha1.ParamValue{
				ObjectVal: map[string]string{"url": "test", "path": "test"},
			},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name: "parametername",
			Type: v1alpha1.ParamTypeObject,
			Default: &v1alpha1.ParamValue{
				ObjectVal: map[string]string{"url": "test", "path": "test"},
			},
		},
	}, {
		name: "inferred type from properties - PropertySpec type is not provided",
		before: &v1alpha1.ParamSpec{
			Name:       "parametername",
			Properties: map[string]v1alpha1.PropertySpec{"key1": {}},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name:       "parametername",
			Type:       v1alpha1.ParamTypeObject,
			Properties: map[string]v1alpha1.PropertySpec{"key1": {Type: "string"}},
		},
	}, {
		name: "inferred type from properties - PropertySpec type is provided",
		before: &v1alpha1.ParamSpec{
			Name:       "parametername",
			Properties: map[string]v1alpha1.PropertySpec{"key2": {Type: "string"}},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name:       "parametername",
			Type:       v1alpha1.ParamTypeObject,
			Properties: map[string]v1alpha1.PropertySpec{"key2": {Type: "string"}},
		},
	}, {
		name: "fully defined ParamSpec - array",
		before: &v1alpha1.ParamSpec{
			Name:        "parametername",
			Type:        v1alpha1.ParamTypeArray,
			Description: "a description",
			Default: &v1alpha1.ParamValue{
				ArrayVal: []string{"array"},
			},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name:        "parametername",
			Type:        v1alpha1.ParamTypeArray,
			Description: "a description",
			Default: &v1alpha1.ParamValue{
				ArrayVal: []string{"array"},
			},
		},
	}, {
		name: "fully defined ParamSpec - object",
		before: &v1alpha1.ParamSpec{
			Name:        "parametername",
			Type:        v1alpha1.ParamTypeObject,
			Description: "a description",
			Default: &v1alpha1.ParamValue{
				ObjectVal: map[string]string{"url": "test", "path": "test"},
			},
		},
		defaultsApplied: &v1alpha1.ParamSpec{
			Name:        "parametername",
			Type:        v1alpha1.ParamTypeObject,
			Description: "a description",
			Default: &v1alpha1.ParamValue{
				ObjectVal: map[string]string{"url": "test", "path": "test"},
			},
		},
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			tc.before.SetDefaults(ctx)
			if d := cmp.Diff(tc.defaultsApplied, tc.before); d != "" {
				t.Error(ktesting.PrintDiffWantGot(d))
			}
		})
	}
}

type ParamValuesHolder struct {
	AOrS v1alpha1.ParamValue `json:"val"`
}

func TestParamValues_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  map[string]interface{}
		result v1alpha1.ParamValue
	}{
		{
			input:  map[string]interface{}{"val": 123},
			result: *v1alpha1.NewStructuredValues("123"),
		},
		{
			input:  map[string]interface{}{"val": "123"},
			result: *v1alpha1.NewStructuredValues("123"),
		},
		{
			input:  map[string]interface{}{"val": ""},
			result: *v1alpha1.NewStructuredValues(""),
		},
		{
			input:  map[string]interface{}{"val": nil},
			result: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeString, ArrayVal: nil},
		},
		{
			input:  map[string]interface{}{"val": []string{}},
			result: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{}},
		},
		{
			input:  map[string]interface{}{"val": []string{"oneelement"}},
			result: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{"oneelement"}},
		},
		{
			input:  map[string]interface{}{"val": []string{"multiple", "elements"}},
			result: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{"multiple", "elements"}},
		},
		{
			input:  map[string]interface{}{"val": map[string]string{"key1": "val1", "key2": "val2"}},
			result: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeObject, ObjectVal: map[string]string{"key1": "val1", "key2": "val2"}},
		},
	}

	for _, c := range cases {
		for _, opts := range []func(enc *json.Encoder){
			// Default encoding
			func(enc *json.Encoder) {},
			// Multiline encoding
			func(enc *json.Encoder) { enc.SetIndent("", "  ") },
		} {
			b := new(bytes.Buffer)
			enc := json.NewEncoder(b)
			opts(enc)
			if err := enc.Encode(c.input); err != nil {
				t.Fatalf("error encoding json: %v", err)
			}

			var result ParamValuesHolder
			if err := json.Unmarshal(b.Bytes(), &result); err != nil {
				t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
			}
			if !reflect.DeepEqual(result.AOrS, c.result) {
				t.Errorf("expected %+v, got %+v", c.result, result)
			}
		}
	}
}

func TestParamValues_UnmarshalJSON_Directly(t *testing.T) {
	cases := []struct {
		desc     string
		input    string
		expected v1alpha1.ParamValue
	}{
		{desc: "empty value", input: ``, expected: *v1alpha1.NewStructuredValues("")},
		{desc: "int value", input: `1`, expected: *v1alpha1.NewStructuredValues("1")},
		{desc: "int array", input: `[1,2,3]`, expected: *v1alpha1.NewStructuredValues("[1,2,3]")},
		{desc: "nested array", input: `[1,\"2\",3]`, expected: *v1alpha1.NewStructuredValues(`[1,\"2\",3]`)},
		{desc: "string value", input: `hello`, expected: *v1alpha1.NewStructuredValues("hello")},
		{desc: "array value", input: `["hello","world"]`, expected: *v1alpha1.NewStructuredValues("hello", "world")},
		{desc: "object value", input: `{"hello":"world"}`, expected: *v1alpha1.NewObject(map[string]string{"hello": "world"})},
	}

	for _, c := range cases {
		v := v1alpha1.ParamValue{}
		if err := v.UnmarshalJSON([]byte(c.input)); err != nil {
			t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
		}
		if !reflect.DeepEqual(v, c.expected) {
			t.Errorf("Failed to unmarshal input '%v': expected %+v, got %+v", c.input, c.expected, v)
		}
	}
}

func TestParamValues_UnmarshalJSON_Error(t *testing.T) {
	cases := []struct {
		desc  string
		input string
	}{
		{desc: "empty value", input: "{\"val\": }"},
		{desc: "wrong beginning value", input: "{\"val\": @}"},
	}

	for _, c := range cases {
		var result ParamValuesHolder
		if err := json.Unmarshal([]byte(c.input), &result); err == nil {
			t.Errorf("Should return err but got nil '%v'", c.input)
		}
	}
}

func TestParamValues_MarshalJSON(t *testing.T) {
	cases := []struct {
		input  v1alpha1.ParamValue
		result string
	}{
		{*v1alpha1.NewStructuredValues("123"), "{\"val\":\"123\"}"},
		{*v1alpha1.NewStructuredValues("123", "1234"), "{\"val\":[\"123\",\"1234\"]}"},
		{*v1alpha1.NewStructuredValues("a", "a", "a"), "{\"val\":[\"a\",\"a\",\"a\"]}"},
		{*v1alpha1.NewObject(map[string]string{"key1": "var1", "key2": "var2"}), "{\"val\":{\"key1\":\"var1\",\"key2\":\"var2\"}}"},
	}

	for _, c := range cases {
		input := ParamValuesHolder{c.input}
		result, err := json.Marshal(&input)
		if err != nil {
			t.Errorf("Failed to marshal input '%v': %v", input, err)
		}
		if string(result) != c.result {
			t.Errorf("Failed to marshal input '%v': expected: %+v, got %q", input, c.result, string(result))
		}
	}
}

func TestExtractNames(t *testing.T) {
	tests := []struct {
		name   string
		params v1alpha1.Params
		want   sets.Set[string]
	}{{
		name:   "no params",
		params: v1alpha1.Params{{}},
		want:   sets.New(""),
	}, {
		name: "extract param names from ParamTypeString",
		params: v1alpha1.Params{{
			Name: "IMAGE", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeString, StringVal: "image-1"},
		}, {
			Name: "DOCKERFILE", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeString, StringVal: "path/to/Dockerfile1"},
		}},
		want: sets.New("IMAGE", "DOCKERFILE"),
	}, {
		name: "extract param names from ParamTypeArray",
		params: v1alpha1.Params{{
			Name: "GOARCH", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{"linux/amd64", "linux/ppc64le", "linux/s390x"}},
		}},
		want: sets.New("GOARCH"),
	}, {
		name: "extract param names from ParamTypeString and ParamTypeArray",
		params: v1alpha1.Params{{
			Name: "GOARCH", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{"linux/amd64", "linux/ppc64le", "linux/s390x"}},
		}, {
			Name: "IMAGE", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeString, StringVal: "image-1"},
		}},
		want: sets.New("GOARCH", "IMAGE"),
	}, {
		name: "extract param name from duplicate params",
		params: v1alpha1.Params{{
			Name: "duplicate", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeArray, ArrayVal: []string{"linux/amd64", "linux/ppc64le", "linux/s390x"}},
		}, {
			Name: "duplicate", Value: v1alpha1.ParamValue{Type: v1alpha1.ParamTypeString, StringVal: "image-1"},
		}},
		want: sets.New("duplicate"),
	}}
	for _, tt := range tests {
		if d := cmp.Diff(tt.want, v1alpha1.Params.ExtractNames(tt.params)); d != "" {
			t.Error(ktesting.PrintDiffWantGot(d))
		}
	}
}

func TestGetNames(t *testing.T) {
	tcs := []struct {
		name   string
		params v1alpha1.ParamSpecs
		want   []string
	}{{
		name: "names from param spec",
		params: v1alpha1.ParamSpecs{{
			Name: "foo",
		}, {
			Name: "bar",
		}},
		want: []string{"foo", "bar"},
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.params.GetNames()
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Error(ktesting.PrintDiffWantGot(d))
			}
		})
	}
}

func TestSortByType(t *testing.T) {
	tcs := []struct {
		name   string
		params v1alpha1.ParamSpecs
		want   []v1alpha1.ParamSpecs
	}{{
		name: "sort by type",
		params: v1alpha1.ParamSpecs{{
			Name: "array1",
			Type: "array",
		}, {
			Name: "string1",
			Type: "string",
		}, {
			Name: "object1",
			Type: "object",
		}, {
			Name: "array2",
			Type: "array",
		}, {
			Name: "string2",
			Type: "string",
		}, {
			Name: "object2",
			Type: "object",
		}},
		want: []v1alpha1.ParamSpecs{
			{{
				Name: "string1",
				Type: "string",
			}, {
				Name: "string2",
				Type: "string",
			}},
			{{
				Name: "array1",
				Type: "array",
			}, {
				Name: "array2",
				Type: "array",
			}},
			{{
				Name: "object1",
				Type: "object",
			}, {
				Name: "object2",
				Type: "object",
			}},
		},
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s, a, o := tc.params.SortByType()
			got := []v1alpha1.ParamSpecs{s, a, o}
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Error(ktesting.PrintDiffWantGot(d))
			}
		})
	}
}

func TestValidateNoDuplicateNames(t *testing.T) {
	tcs := []struct {
		name          string
		params        v1alpha1.ParamSpecs
		expectedError *apis.FieldError
	}{{
		name: "no duplicates",
		params: v1alpha1.ParamSpecs{{
			Name: "foo",
		}, {
			Name: "bar",
		}},
	}, {
		name: "duplicates",
		params: v1alpha1.ParamSpecs{{
			Name: "foo",
		}, {
			Name: "foo",
		}},
		expectedError: &apis.FieldError{
			Message: `parameter appears more than once`,
			Paths:   []string{"params[foo]"},
		},
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.params.ValidateNoDuplicateNames()
			if d := cmp.Diff(tc.expectedError.Error(), got.Error()); d != "" {
				t.Error(ktesting.PrintDiffWantGot(d))
			}
		})
	}
}
