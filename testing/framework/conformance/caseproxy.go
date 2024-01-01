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
	"slices"
	"strings"

	"github.com/katanomi/pkg/testing/framework/base"
	. "github.com/onsi/ginkgo/v2"
)

// NewCaseProxy construct a new proxy for a special testcase
func NewCaseProxy(testCase *testCase) *caseProxy {
	proxy := &caseProxy{
		testCase:  testCase,
		proxyNode: testCase.Node.Clone(),
	}
	return proxy
}

// caseProxy is a proxy for testcase to connect testcase to higher-level nodes
type caseProxy struct {
	*testCase

	proxyNode         *Node
	focusedTestPoints []*testPoint
	caseRegister      func(ctx context.Context)
}

// LinkParentNode link test case to parent node
func (m *caseProxy) LinkParentNode(node *Node) {
	m.proxyNode.LinkParentNode(node)
	m.proxyNode.caseRegister = m.RegisterTestCase
}

// RegisterTestCase register all the test cases
func (m *caseProxy) RegisterTestCase() {
	if m.caseRegister == nil {
		return
	}

	var labels Labels
	if parentNode := m.proxyNode.ParentNode; parentNode != nil {
		labels = append(labels, strings.Join(parentNode.Labels(), "#"))
	}
	ctx := base.WithContextLabel(context.Background(), labels)

	m.caseRegister(ctx)
}

// New construct a new case set makes it possible to isolate different scenarios
func (m *caseProxy) New() CaseSet {
	clone := &caseProxy{
		testCase:          m.testCase,
		proxyNode:         m.proxyNode.Clone(),
		focusedTestPoints: []*testPoint{},
		caseRegister:      m.caseRegister,
	}
	return clone
}

// Build convert to CaseSet interface
func (m *caseProxy) Build(caseRegister func(ctx context.Context)) CaseSetFactory {
	m.caseRegister = caseRegister
	return m
}

// Focus specify the test points to be executed
func (m *caseProxy) Focus(testPoints ...*testPoint) CaseSet {
	for _, item := range testPoints {
		if slices.Contains(m.testPoints, item) {
			m.focusedTestPoints = append(m.focusedTestPoints, testPoints...)
		}
	}
	var subNodes []*Node
	for _, item := range m.testCase.testPoints {
		if slices.Contains(m.focusedTestPoints, item) {
			subNodes = append(subNodes, item.Node.Clone())
		}
	}

	m.proxyNode.SubNodes = subNodes
	return m
}
