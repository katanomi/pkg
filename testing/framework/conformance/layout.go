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

package conformance

func NewModuleCase(moduleName string) *ModuleCase {
	m := &ModuleCase{
		Node: NewNode(ModuleLevel, moduleName),
	}
	pushModuleCase(m)
	return m
}

type ModuleCase struct {
	*Node
}

func (m *ModuleCase) AddFeatureCase(features ...*featureCase) {
	for _, feature := range features {
		feature.LinkParentNode(m.Node)
	}
}

func NewFeatureCase(featureName string, caseSets ...CaseSet) *featureCase {
	featureCases := &featureCase{
		Node: NewNode(FeatureLevel, featureName),
	}

	for _, caseSet := range caseSets {
		featureCases.AddTestCaseSet(caseSet)
	}

	return featureCases
}

type featureCase struct {
	*Node
}

func (f *featureCase) AddTestCaseSet(caseSets ...CaseSet) *featureCase {
	for _, caseSet := range caseSets {
		caseSet.LinkParentNode(f.Node)
	}
	return f
}
