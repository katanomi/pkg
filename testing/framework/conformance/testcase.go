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
)

func NewTestCase(name string) *testCase {
	return &testCase{node: NewNode(TestCaseLevel, name)}
}

type testCase struct {
	node *Node

	testPoints []*testPoint
}

func (t *testCase) Build(caseRegister func(ctx context.Context)) CaseSetFactory {
	return NewCaseProxy(t).Build(caseRegister)
}

// NewTestPoint create a test point
func (t *testCase) NewTestPoint(name string) *testPoint {
	tp := NewTestPoint(name)
	tp.node.LinkParentNode(t.node)

	t.testPoints = append(t.testPoints, tp)
	return tp
}
