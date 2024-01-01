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

// CaseSetFactory factory to create a case set
type CaseSetFactory interface {
	// New construct a new case set
	New() CaseSet
}

// CaseSet describe a set of test cases
type CaseSet interface {
	// Focus specify the test points to be executed
	Focus(testPoints ...*testPoint) CaseSet

	// LinkParentNode link test case to parent node
	LinkParentNode(node *Node)
}

// CustomAssertion describe a interface to add custom assertion for a special testpoint
type CustomAssertion interface {
	// AddAssertion add a custom assertion
	AddAssertion(f interface{}) *testPoint
}

// AddAssertionFunc helper function to make a implementation of CustomAssertion interface
type AddAssertionFunc func(f interface{}) *testPoint

// AddAssertion add a custom assertion
func (p AddAssertionFunc) AddAssertion(f interface{}) *testPoint {
	return p(f)
}
