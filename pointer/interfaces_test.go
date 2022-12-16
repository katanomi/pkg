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
package pointer

import (
	"testing"
)

func TestIsNil(T *testing.T) {
	type test struct {
		name   string
		input  interface{}
		expect bool
	}

	type testWarp struct {
		input *test
	}
	tests := []test{
		{"nil", nil, true},
		{"nil string", "nil", false},

		{"test{}", test{}, false},
		{"&test{}", &test{}, false},

		{"test{input:nil}", test{input: nil}, false},
		{"&test{input:nil}", &test{input: nil}, false},

		{"testWarp", testWarp{nil}, false},
		{"&testWarp", &testWarp{nil}, false},
	}
	for _, t := range tests {
		if r := IsNil(t.input); r != t.expect {
			T.Errorf("expect %v but get %v for %s", t.expect, r, t.name)
		}
	}
}
