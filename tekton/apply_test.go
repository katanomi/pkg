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

package tekton

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func Test_IsClustertemplate(t *testing.T) {
	RegisterTestingT(t)

	expectObjects := map[string]map[string]string{
		"params.Object": {
			"object1": "{abc:def}",
		},
		"params[\"Object\"]": {
			"object1": "{abc:def}",
		},
		"params['Object']": {
			"object1": "{abc:def}",
		},
	}
	expectArrays := map[string][]string{
		"params[\"Array\"]": {"array1", "array2"},
		"params['Array']":   {"array1", "array2"},
		"params.Array":      {"array1", "array2"},
	}
	expectStrings := map[string]string{
		"params.String":      "string",
		"params[\"String\"]": "string",
		"params['String']":   "string",

		"params[\"Array\"][0]": "array1",
		"params.Array[0]":      "array1",
		"params['Array'][0]":   "array1",

		"params.Array[1]":      "array2",
		"params['Array'][1]":   "array2",
		"params[\"Array\"][1]": "array2",

		"params.Object.object1": "{abc:def}",
	}

	for _, c := range []struct {
		description   string
		ctx           context.Context
		paramSpec     []v1beta1.ParamSpec
		params        []v1beta1.Param
		expectStrings map[string]string
		expectArrays  map[string][]string
		expectObjects map[string]map[string]string
	}{
		{
			description:   "empty",
			ctx:           context.Background(),
			paramSpec:     []v1beta1.ParamSpec{},
			params:        []v1beta1.Param{},
			expectObjects: map[string]map[string]string{},
			expectArrays:  map[string][]string{},
			expectStrings: map[string]string{},
		},
		{
			description: "use default strings arrays and objects",
			ctx:         context.Background(),
			paramSpec: []v1beta1.ParamSpec{
				{
					Name: "String",
					Type: v1beta1.ParamTypeString,
					Default: &v1beta1.ParamValue{
						Type:      v1beta1.ParamTypeString,
						StringVal: "string",
					},
				},
				{
					Name: "Array",
					Type: v1beta1.ParamTypeArray,
					Default: &v1beta1.ParamValue{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: []string{"array1", "array2"},
					},
				},
				{
					Name: "Object",
					Type: v1beta1.ParamTypeArray,
					Default: &v1beta1.ParamValue{
						Type:      v1beta1.ParamTypeObject,
						ObjectVal: map[string]string{"object1": "{abc:def}"},
					},
				},
			},
			params:        []v1beta1.Param{},
			expectObjects: expectObjects,
			expectArrays:  expectArrays,
			expectStrings: expectStrings,
		},

		{
			description: "use default and value strings arrays and objects value will override default",
			ctx:         context.Background(),
			paramSpec: []v1beta1.ParamSpec{
				{
					Name: "String",
					Type: v1beta1.ParamTypeString,
					Default: &v1beta1.ParamValue{
						Type:      v1beta1.ParamTypeString,
						StringVal: "default",
					},
				},
				{
					Name: "Object",
					Type: v1beta1.ParamTypeArray,
					Default: &v1beta1.ParamValue{
						Type:      v1beta1.ParamTypeObject,
						ObjectVal: map[string]string{"object1": "{abc:def}"},
					},
				},
			},
			params: []v1beta1.Param{
				{
					Name: "String",
					Value: v1beta1.ParamValue{
						Type:      v1beta1.ParamTypeString,
						StringVal: "string",
					},
				},
				{
					Name: "Array",
					Value: v1beta1.ParamValue{
						Type:     v1beta1.ParamTypeArray,
						ArrayVal: []string{"array1", "array2"},
					},
				},
			},
			expectObjects: expectObjects,
			expectArrays:  expectArrays,
			expectStrings: expectStrings,
		},
	} {

		t.Logf("<=== starting %s...", c.description)
		strings, arrays, objects := Replacements(c.ctx, c.paramSpec, c.params)
		Expect(strings).To(Equal(c.expectStrings))
		Expect(arrays).To(Equal(c.expectArrays))
		Expect(objects).To(Equal(c.expectObjects))
		t.Logf("===> passed %s...", c.description)
	}
}
