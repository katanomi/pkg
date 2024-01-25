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
	"fmt"

	"github.com/tektoncd/pipeline/pkg/apis/config"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var (
	paramPatterns = []string{
		"params.%s",
		"params[%q]",
		"params['%s']",
	}
)

const (
	// objectIndividualVariablePattern is the reference pattern for object individual keys params.<object_param_name>.<key_name>
	objectIndividualVariablePattern = "params.%s.%s"
)

// reference from  https://github.com/katanomi/pipeline/blob/a9210847e1cb183797379917bec8dfd221450321/pkg/reconciler/pipelinerun/resources/apply.go#L106
// and do some refactor base that

// Replacements return replacements base on the params spec and provided params values
func Replacements(ctx context.Context, paramSpecs []v1beta1.ParamSpec, params []v1beta1.Param) (stringReplacements map[string]string, arrayReplacements map[string][]string, objectReplacements map[string]map[string]string) {

	ctx = WithDefaultConfig(ctx)

	strings, arrays, objects := paramDefaultReplacements(ctx, paramSpecs)
	// Set and overwrite params with the ones from the parameters provided
	valueStrings, valueArrays, valueObjects := paramValueReplacements(ctx, params)

	for k, v := range valueStrings {
		strings[k] = v
	}
	for k, v := range valueArrays {
		arrays[k] = v
	}
	for k, v := range valueObjects {
		objects[k] = v
	}

	return strings, arrays, objects
}

func paramDefaultReplacements(ctx context.Context, paramSpecs []v1beta1.ParamSpec) (map[string]string, map[string][]string, map[string]map[string]string) {

	cfg := config.FromContextOrDefaults(ctx)

	stringReplacements := map[string]string{}
	arrayReplacements := map[string][]string{}
	objectReplacements := map[string]map[string]string{}

	// Set all the default replacements
	for _, p := range paramSpecs {
		if p.Default == nil {
			continue
		}
		switch p.Default.Type {
		case v1beta1.ParamTypeArray:
			for _, pattern := range paramPatterns {
				// array indexing for param is alpha feature
				if cfg.FeatureFlags.EnableAPIFields == config.AlphaAPIFields {
					for i := 0; i < len(p.Default.ArrayVal); i++ {
						stringReplacements[fmt.Sprintf(pattern+"[%d]", p.Name, i)] = p.Default.ArrayVal[i]
					}
				}
				arrayReplacements[fmt.Sprintf(pattern, p.Name)] = p.Default.ArrayVal
			}
		case v1beta1.ParamTypeObject:
			for _, pattern := range paramPatterns {
				objectReplacements[fmt.Sprintf(pattern, p.Name)] = p.Default.ObjectVal
			}
			for k, v := range p.Default.ObjectVal {
				stringReplacements[fmt.Sprintf(objectIndividualVariablePattern, p.Name, k)] = v
			}
		default:
			for _, pattern := range paramPatterns {
				stringReplacements[fmt.Sprintf(pattern, p.Name)] = p.Default.StringVal
			}
		}
	}
	return stringReplacements, arrayReplacements, objectReplacements
}

func paramValueReplacements(ctx context.Context, params []v1beta1.Param) (map[string]string, map[string][]string, map[string]map[string]string) {
	// stringReplacements is used for standard single-string stringReplacements,
	// while arrayReplacements/objectReplacements contains arrays/objects that need to be further processed.
	stringReplacements := map[string]string{}
	arrayReplacements := map[string][]string{}
	objectReplacements := map[string]map[string]string{}
	cfg := config.FromContextOrDefaults(ctx)

	for _, p := range params {
		switch p.Value.Type {
		case v1beta1.ParamTypeArray:
			for _, pattern := range paramPatterns {
				// array indexing for param is alpha feature
				if cfg.FeatureFlags.EnableAPIFields == config.AlphaAPIFields {
					for i := 0; i < len(p.Value.ArrayVal); i++ {
						stringReplacements[fmt.Sprintf(pattern+"[%d]", p.Name, i)] = p.Value.ArrayVal[i]
					}
				}
				arrayReplacements[fmt.Sprintf(pattern, p.Name)] = p.Value.ArrayVal
			}
		case v1beta1.ParamTypeObject:
			for _, pattern := range paramPatterns {
				objectReplacements[fmt.Sprintf(pattern, p.Name)] = p.Value.ObjectVal
			}
			for k, v := range p.Value.ObjectVal {
				stringReplacements[fmt.Sprintf(objectIndividualVariablePattern, p.Name, k)] = v
			}
		default:
			for _, pattern := range paramPatterns {
				stringReplacements[fmt.Sprintf(pattern, p.Name)] = p.Value.StringVal
			}
		}
	}

	return stringReplacements, arrayReplacements, objectReplacements
}
