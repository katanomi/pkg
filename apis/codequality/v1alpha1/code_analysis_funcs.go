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
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (AnalysisResult) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *AnalysisResult) {
	if values != nil {
		result = &AnalysisResult{
			Result:    values[path.Child("result").String()],
			ReportURL: values[path.Child("reportURL").String()],
			TaskID:    values[path.Child("taskID").String()],
			ProjectID: values[path.Child("projectID").String()],
			Metrics:   AnalisysMetrics{}.GetObjectWithValues(ctx, path.Child("metrics"), values),
		}
	}
	return
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (AnalisysMetrics) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *AnalisysMetrics) {
	if values != nil {
		result = &AnalisysMetrics{
			Branch:    CodeChangeMetrics{}.GetObjectWithValues(ctx, path.Child("branch"), values),
			Target:    CodeChangeMetrics{}.GetObjectWithValues(ctx, path.Child("target"), values),
			Languages: strings.Split(values[path.Child("languages").String()], ","),
			Ratings:   map[string]AnalysisRating{},
			CodeSize:  CodeSize{}.GetObjectWithValues(ctx, path.Child("codeSize"), values),
		}

		// because it is a map the only way is to filter out
		// possible keys then assign value
		ratingsKey := path.Child("ratings")
		ratingsPrefix := ratingsKey.String()

		for k := range values {
			if strings.HasPrefix(k, ratingsPrefix) {
				// remove the prefix and suffix
				// "abc.ratings." / "vunerability" / "rate"
				// taking the middle part
				ratingType := strings.SplitN(strings.TrimPrefix(k, ratingsPrefix+"."), ".", 2)[0]
				if rating := result.Ratings[""].GetObjectWithValues(ctx, ratingsKey.Child(ratingType), values); rating != nil {
					result.Ratings[ratingType] = *rating
				}
			}
		}
	}
	return
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (CodeChangeMetrics) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *CodeChangeMetrics) {
	if values != nil {
		result = &CodeChangeMetrics{}
		if coverage := result.CoverageRate.GetObjectWithValues(ctx, path.Child("coverage"), values); coverage != nil {
			result.CoverageRate = *coverage
		}

		if duplications := result.DuplicationRate.GetObjectWithValues(ctx, path.Child("duplications"), values); duplications != nil {
			result.DuplicationRate = *duplications
		}
	}
	return

}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (CodeChangeRates) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *CodeChangeRates) {
	if values != nil {
		result = &CodeChangeRates{
			New:   values[path.Child("new").String()],
			Total: values[path.Child("total").String()],
		}
	}
	return
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (AnalysisRating) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *AnalysisRating) {
	if values != nil {
		result = &AnalysisRating{
			Rate:        values[path.Child("rate").String()],
			IssuesCount: strconvAtoi(values[path.Child("issues").String()]),
		}
	}
	return
}

// GetObjectWithValues inits an object based on a json.path values map
// returns nil if values is nil
func (CodeSize) GetObjectWithValues(ctx context.Context, path *field.Path, values map[string]string) (result *CodeSize) {
	if values != nil {
		result = &CodeSize{}
		result.LinesOfCode, _ = strconv.Atoi(values[path.Child("linesOfCode").String()])
	}
	return
}
