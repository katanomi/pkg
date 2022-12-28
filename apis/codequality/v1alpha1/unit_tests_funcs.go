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

package v1alpha1

import (
	"context"
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// IsEmpty returns true if the struct is empty
func (u UnitTestsResult) IsEmpty() bool {
	return u.Coverage.IsEmpty() && u.TestResult.IsEmpty()
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (UnitTestsResult) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *UnitTestsResult) {
	if values != nil {
		result = &UnitTestsResult{
			TestResult: TestResult{}.GetObjectWithValues(ctx, path.Child("testResults"), values),
			Coverage:   TestCoverage{}.GetObjectWithValues(ctx, path.Child("coverage"), values),
		}
	}
	return
}

// IsEmpty returns true if the struct is empty
func (c *TestCoverage) IsEmpty() bool {
	if c == nil {
		return true
	}
	return c.Branches == "" && c.Lines == ""
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (TestCoverage) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *TestCoverage) {
	if values != nil {
		result = &TestCoverage{
			Lines:    values[path.Child("lines").String()],
			Branches: values[path.Child("branches").String()],
		}
	}
	return
}

// IsEmpty returns true if the struct is empty
func (t *TestResult) IsEmpty() bool {
	if t == nil {
		return true
	}
	return t.Failed+t.Passed+t.Skipped == 0 && t.PassedTestsRate == ""
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (TestResult) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *TestResult) {
	if values != nil {
		result = &TestResult{
			Passed:          strconvAtoi(values[path.Child("passed").String()]),
			Failed:          strconvAtoi(values[path.Child("failed").String()]),
			Skipped:         strconvAtoi(values[path.Child("skipped").String()]),
			PassedTestsRate: values[path.Child("passedTestsRate").String()],
		}

		if errNumStr, errExists := values[path.Child("error").String()]; errExists {
			errNum := strconvAtoi(errNumStr)
			result.Failed += errNum
		}
	}
	return
}

// PassedTestsRate rate of passed tests calculated using passed / (passed + failed) * 100, i.e 96.54
// if input t is nil, will return "0"
func PassedTestsRate(t *TestResult) string {
	if t == nil {
		return "0"
	}

	passedTestsRate := float64(t.Passed) / float64(t.Passed+t.Failed)
	passedTestsRate *= 100
	return fmt.Sprintf("%.2f", passedTestsRate)
}

func strconvAtoi(stringVal string) (val int) {
	val, _ = strconv.Atoi(stringVal)
	return
}
