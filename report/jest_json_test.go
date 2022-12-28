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
	"reflect"
	"testing"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
)

func TestResult_JestJson(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		wantTestResult v1alpha1.TestResult
		wantErr        bool
	}{
		{
			name: "parser success",
			path: "./testdata/jest_json_success.json",
			wantTestResult: v1alpha1.TestResult{
				Passed:          5,
				Skipped:         2,
				Failed:          2,
				PassedTestsRate: "71.43",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JestJsonParser{}
			testResult, err := p.Parse(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JestJsonParser.TestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			converter, _ := testResult.(ConvertToTestResult)
			outResult := converter.ConvertToTestResult()
			if !reflect.DeepEqual(outResult, tt.wantTestResult) {
				t.Errorf("JestJsonParser.TestResult() = %v, want %v", outResult, tt.wantTestResult)
			}
		})
	}
}
