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

type testPoint struct {
	*Node

	testCase *testCase

	additionalAssertions map[*featureCase]interface{}
}

func (t *testPoint) Labels(ctx context.Context) Labels {
	contextLabels := append(base.ContextLabel(ctx), t.Node.Labels()...)
	return Labels{
		strings.Join(t.Node.Labels(), "#"),
		strings.Join(contextLabels, "#"),
	}
}

func (t *testPoint) CheckExternalAssertion(args ...interface{}) {
	contextLabel := CurrentSpecReport().Labels()
	for feature, assertFunc := range t.additionalAssertions {
		featureLabel := strings.Join(feature.Node.Labels(), "#")
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
func (t *testPoint) Bind(feature *featureCase) CustomAssertion {
	return AddAssertionFunc(func(f interface{}) *testPoint {
		return t.AddAssertion(feature, f)
	})
}

func (t *testPoint) AddAssertion(feature *featureCase, assertFunc interface{}) *testPoint {
	// check assertFunc is a function
	val := reflect.ValueOf(assertFunc)
	if val.Kind() != reflect.Func || val.Type().NumIn() == 0 {
		// TODO
		panic("assertFunc must be a function with at least one argument")
	}
	if t.additionalAssertions == nil {
		t.additionalAssertions = make(map[*featureCase]interface{})
	}

	t.additionalAssertions[feature] = assertFunc
	return t
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
