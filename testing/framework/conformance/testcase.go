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

func NewTestCase(name string) *testCase {
	return &testCase{Node: NewNode(TestCaseLevel, name)}
}

type testCase struct {
	*Node

	testPoints []*testPoint
}

func (t *testCase) Proxy() *caseProxy {
	return NewCaseProxy(t)
}

func (t *testCase) Labels() Labels {
	return t.Node.LevelLabel.Labels()
}

// NewTestPoint create a test point
func (t *testCase) NewTestPoint(name string) *testPoint {
	tp := &testPoint{
		Node:     NewNode(TestPointLevel, name),
		testCase: t,
	}

	tp.LinkParentNode(t.Node)

	t.testPoints = append(t.testPoints, tp)

	return tp
}
