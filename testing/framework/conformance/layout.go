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

import (
	. "github.com/onsi/ginkgo/v2"
)

// currently, only support 3 levels of hierarchy(module -> function -> feature).
// But it can be extended to support more levels of hierarchy easily.

// NewModuleCase construct a new module case
func NewModuleCase(moduleName string) *moduleCase {
	m := &moduleCase{
		node: NewNode(ModuleLevel, moduleName),
	}
	pushModuleCase(m)
	return m
}

type moduleCase struct {
	node *Node
}

// AddFunctionCase register feature case
func (m *moduleCase) AddFunctionCase(functionCases ...*functionCase) {
	for _, fCase := range functionCases {
		fCase.node.LinkParentNode(m.node)
	}
}

// AddFeatureCase register feature case
// it will create a virtual function case to hold feature cases,
// and the name of virtual function case is the same as module case
func (m *moduleCase) AddFeatureCase(featureCases ...*featureCase) {
	virtualFunctionCase := NewFunctionCase(m.node.Name, featureCases...)
	m.AddFunctionCase(virtualFunctionCase)
}

func (m *moduleCase) RegisterTestCase() {
	m.node.RegisterTestCase()
}

// NewFunctionCase construct a new feature case
func NewFunctionCase(functionName string, featureCases ...*featureCase) *functionCase {
	fCase := &functionCase{
		node: NewNode(FunctionLevel, functionName),
	}

	fCase.AddFeature(featureCases...)
	return fCase
}

type functionCase struct {
	node *Node
}

// AddFeature register subFeature case to the feature case
func (f *functionCase) AddFeature(features ...*featureCase) *functionCase {
	for _, feature := range features {
		feature.node.LinkParentNode(f.node)
	}
	return f
}

// NewFeatureCase construct a new subFeature case
func NewFeatureCase(featureName string, caseSets ...CaseSet) *featureCase {
	fCase := &featureCase{
		node: NewNode(FeatureLevel, featureName),
	}

	fCase.AddTestCaseSet(caseSets...)
	return fCase
}

// featureCase represents a test case for a specific feature.
// It contains a node that represents the feature in the test hierarchy.
type featureCase struct {
	node *Node
}

// AddTestCaseSet register testcase to the subFeature case
func (f *featureCase) AddTestCaseSet(caseSets ...CaseSet) *featureCase {
	for _, caseSet := range caseSets {
		caseSet.LinkParentNode(f.node)
		caseSet.bindFeature(f)
	}
	return f
}

// Labels returns the labels associated with the featureCase.
func (f *featureCase) Labels() Labels {
	return f.node.Labels()
}

// lazyFeatureCaseBind is a struct used to lazily bind custom assertion to a feature case
type lazyFeatureCaseBind struct {
	feature    *featureCase
	assertFunc interface{}
}

// Labels returns the labels associated with the featureCase.
// If there is no associated feature case, it will return an empty slice of labels.
func (c *lazyFeatureCaseBind) Labels() Labels {
	if c.feature == nil {
		return Labels{}
	}

	return c.feature.Labels()
}
