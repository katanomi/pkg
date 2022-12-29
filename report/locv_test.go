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

func TestLcovParser_Coverage(t *testing.T) {
	tests := map[string]struct {
		path             string
		wantTestCoverage v1alpha1.TestCoverage
		wantErr          error
	}{
		"parse lcov success": {
			path: "./testdata/lcovparser-success.info",
			wantTestCoverage: v1alpha1.TestCoverage{
				Lines:    "70.00",
				Branches: "50.00",
			},
		},
		"lcov file not found": {
			path:    "./testdata/lcovparser-not-found.info",
			wantErr: fmt.Errorf("open ./testdata/lcovparser-not-found.info: no such file or directory"),
		},
		"lcov parseline failed": {
			path:    "./testdata/lcovparser-failed.info",
			wantErr: fmt.Errorf(`invalid lcov text:LF:f. error: strconv.Atoi: parsing "f": invalid syntax`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &LcovParser{}
			testConverager, err := p.Parse(tt.path)
			if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
				t.Errorf("LcovParser.TestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if testConverager == nil {
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
