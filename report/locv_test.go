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

func TestLcovParser_Coverage(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		wantTestCoverage v1alpha1.TestCoverage
		wantErr          bool
	}{
		{
			name: "paser success.",
			path: "./testdata/lcovparser-success.info",
			wantTestCoverage: v1alpha1.TestCoverage{
				Lines:    "70.00",
				Branches: "50.00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &LcovParser{}
			testConverager, err := p.Parse(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LcovParser.TestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			converter, _ := testConverager.(ConvertToTestCoverage)
			outResult := converter.ConvertToTestCoverage()
			if !reflect.DeepEqual(outResult, tt.wantTestCoverage) {
				t.Errorf("LcovParser.TestResult() = %v, want %v", outResult, tt.wantTestCoverage)
			}
		})
	}
}
