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

package report

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
)

func TestResult_MochaJson(t *testing.T) {
	tests := map[string]struct {
		path           string
		wantTestResult v1alpha1.TestResult
		wantErr        error
	}{
		"parse mocha success": {
			path: "./testdata/mochajsonparser-success.json",
			wantTestResult: v1alpha1.TestResult{
				Passed:          3,
				Skipped:         2,
				Failed:          1,
				PassedTestsRate: "75.00",
			},
		},
		"parse file not found": {
			path:    "./testdata/mochajsonparser-not-found.json",
			wantErr: fmt.Errorf("open ./testdata/mochajsonparser-not-found.json: no such file or directory"),
		},
		"parse unmarshal failed": {
			path:    "./testdata/mochajsonparser-failed.json",
			wantErr: fmt.Errorf("json: cannot unmarshal string into Go struct field MochaJsonStats.stats.tests of type int"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &MochaJsonParser{}
			testResult, err := p.Parse(tt.path)
			if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
				t.Errorf("MochaJsonParser.TestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if testResult == nil {
				return
			}

			converter, _ := testResult.(ConvertToTestResult)
			outResult := converter.ConvertToTestResult()
			if !reflect.DeepEqual(outResult, tt.wantTestResult) {
				t.Errorf("MochaJsonParser.TestResult() = %v, want %v", outResult, tt.wantTestResult)
			}
		})
	}
}
