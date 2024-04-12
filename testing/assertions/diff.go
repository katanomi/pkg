/*
Copyright 2024 The Katanomi Authors.

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

package assertions

import (
	"fmt"

	"github.com/onsi/gomega/format"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DiffEqual will use github.com/google/go-cmp/cmp.Diff to compare
// using a list of cmp.Option to add any additional cmp configuration
// like ignoring types or fields not necessary in the tests.
// This package offers a few default options like IgnoreObjectMetaFields that can
// be used with kubernetes Objects
func DiffEqual(expected interface{}, cmpOptions ...cmp.Option) gomega.OmegaMatcher {
	return &DiffEqualMatcher{Expected: expected, Options: cmpOptions}
}

// DeepEqual is an alias for DiffEqual
// uses github.com/google/go-cmp/cmp.Diff to compare objects
var DeepEqual = DiffEqual

// DiffEqualMatcher will use github.com/google/go-cmp/cmp.Diff to compare objects
// and can use a list of cmp.Option to fine-tune comparison.

type DiffEqualMatcher struct {
	// Expected object
	Expected interface{}
	// Options for cmp.Diff function. There are multiple available options
	// inside the github.com/google/go-cmp/cmp/cmpopts package
	// as well as IgnoreTypeMetaFields and the sort of package
	Options []cmp.Option

	// DiffCleanFunc are deprecated and should be replaced with []cmp.Option
	DiffCleanFunc []func(object interface{}) interface{}

	diff string
}

// Match compares the actual and expected values using cmp.Diff and given Options
// will return true if the values are equivalent, otherwise it will return false
// with any necessary error messages
func (matcher *DiffEqualMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Both actual and expected must not be nil.")
	}

	// this is a temporary compatibility functionality that will be deprecated
	// when DiffCleanFunc is not used anymore
	if len(matcher.DiffCleanFunc) > 0 {
		for _, cleanFunc := range matcher.DiffCleanFunc {
			actual = cleanFunc(actual)
			matcher.Expected = cleanFunc(matcher.Expected)
		}
	}

	diff := cmp.Diff(matcher.Expected, actual, matcher.Options...)
	matcher.diff = diff
	if diff == "" {
		return true, nil
	}

	return false, nil
}

// FailureMessage returns a message to be displayed if the actual value does not
// match the expected value.
func (matcher *DiffEqualMatcher) FailureMessage(_ interface{}) (message string) {
	return format.Message(matcher.diff, "to diff empty", "")
}

// NegatedFailureMessage returns a message to be displayed if the actual value
// matches the expected value but is expected to not.
func (matcher *DiffEqualMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return format.Message(matcher.diff, "not to diff empty", "")
}

// IgnoreObjectMetaFields will ignore the fields in metav1.ObjectMeta
// that are flaky or hard to create expectations for in tests.
// If no fields are passed, the default set fields are used
// as DefaultIgnoreTypeMetaFields variable
func IgnoreObjectMetaFields(fields ...string) cmp.Option {
	if len(fields) == 0 {
		fields = append(fields, DefaultIgnoreObjectMetaFields...)
	}
	return cmpopts.IgnoreFields(metav1.ObjectMeta{}, fields...)
}

// DefaultIgnoreObjectMetaFields common fields to ignore in metav1.ObjectMeta using cmp.Diff method
var DefaultIgnoreObjectMetaFields = []string{"CreationTimestamp", "DeletionTimestamp", "DeletionGracePeriodSeconds", "Finalizers", "UID", "Generation", "ManagedFields", "ResourceVersion", "SelfLink", "UID"}

// IgnoreTypeMeta will ignore the fields in metav1.TypeMeta in a cmp.Diff comparison
var IgnoreTypeMeta = cmpopts.IgnoreTypes(metav1.TypeMeta{})

// IgnoreRawExtension will ignore the fields in runtime.RawExtension in a cmp.Diff comparison
var IgnoreRawExtension = cmpopts.IgnoreTypes(runtime.RawExtension{})
