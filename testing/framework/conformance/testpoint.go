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
	"context"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/utils/strings/slices"

	"github.com/katanomi/pkg/testing/framework/base"
)

// NewTestPoint construct a new test point
func NewTestPoint(name string) *testPoint {
	return &testPoint{
		node: NewNode(TestPointLevel, name),
	}
}

type testPoint struct {
	node *Node

	// additionalAssertions each feature can have a custom assertion
	additionalAssertions map[*featureCase]interface{}

	// lazyFeatureCaseBinder is a lazy binder for custom assertion
	// It is used to lazily bind the feature case to the custom assertion
	lazyFeatureCaseBinder *lazyFeatureCaseBind
}

// Labels returns all the labels for the test point
// contains the layout labels which get from the context at runtime
func (t *testPoint) Labels(ctx context.Context) Labels {
	contextLabels := append(base.ContextLabel(ctx), t.node.Labels()...)
	return Labels{
		strings.Join(t.node.Labels(), "#"),
		strings.Join(contextLabels, "#"),
	}
}

// CheckExternalAssertion check the external assertion
// It only executes the assertions related to the current feature, not all of them.
func (t *testPoint) CheckExternalAssertion(args ...interface{}) {
	contextLabel := CurrentSpecReport().Labels()
	for feature, assertFunc := range t.additionalAssertions {
		featureLabel := strings.Join(feature.Labels(), "#")
		if !slices.Contains(contextLabel, featureLabel) {
			continue
		}
		// call additional assertions
		func() {
			// TODO: A friendly reminder when reflect is wrong.
			//defer func() {
			//if e := recover(); e != nil {
			//	panic(fmt.Sprintf("[%s] custom assertion failed: %v", featureIdentity, e))
			//}
			//}()

			// TODO: support for more flexible parameters
			// Maybe can pass the appropriate parameters according to
			// the signature of the assertFunc function
			invokeFunction(assertFunc, args)
		}()
	}
}

// Bind alias of AddAssertion
// Deprecated: use AddAssertion instead
// retain this method for compatibility
func (t *testPoint) Bind(_ *featureCase) CustomAssertion {
	return t
}

// AddAssertion add a custom assertion to the test point for a special feature.
// When the method is called, it will create a lazyFeatureCaseBinder to bind the
// custom assertion to the feature case later.
//
// Example:
//
//	auth.TestPointOauth2.AddAssertion(func(statusCode int) {
//	  Expect(statusCode >= 200).To(BeTrue())
//	}),
func (t *testPoint) AddAssertion(assertFunc interface{}) *testPoint {
	// check assertFunc is a function
	val := reflect.ValueOf(assertFunc)
	if val.Kind() != reflect.Func || val.Type().NumIn() == 0 {
		panic("assertFunc must be a function with at least one argument")
	}

	if t.lazyFeatureCaseBinder != nil {
		panic("assertion function should finish binding before adding new assertion")
	}

	t.lazyFeatureCaseBinder = &lazyFeatureCaseBind{assertFunc: assertFunc}
	return t
}

// bindFeature bind the feature case to the current test point
// Now it used to bind the custom assertion to the feature case.
func (t *testPoint) bindFeature(feature *featureCase) {
	if t.lazyFeatureCaseBinder == nil {
		return
	}
	if t.additionalAssertions == nil {
		t.additionalAssertions = make(map[*featureCase]interface{})
	}
	t.additionalAssertions[feature] = t.lazyFeatureCaseBinder.assertFunc
	t.lazyFeatureCaseBinder = nil
}

// invokeFunction invokes the function with the given parameters
func invokeFunction(function interface{}, parameters []interface{}) []reflect.Value {
	inValues := make([]reflect.Value, len(parameters))

	funcType := reflect.TypeOf(function)
	limit := funcType.NumIn()
	if funcType.IsVariadic() {
		limit = limit - 1
	}

	for i := 0; i < limit && i < len(parameters); i++ {
		inValues[i] = computeValue(parameters[i], funcType.In(i))
	}

	if funcType.IsVariadic() {
		variadicType := funcType.In(limit).Elem()
		for i := limit; i < len(parameters); i++ {
			inValues[i] = computeValue(parameters[i], variadicType)
		}
	}

	return reflect.ValueOf(function).Call(inValues)
}

func computeValue(parameter interface{}, t reflect.Type) reflect.Value {
	if parameter == nil {
		return reflect.Zero(t)
	} else {
		return reflect.ValueOf(parameter)
	}
}
