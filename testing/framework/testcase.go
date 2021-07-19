/*
Copyright 2021 The Katanomi Authors.

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

package framework

import (
	"fmt"

	"github.com/onsi/ginkgo"
)

// TestCasePriority priority for the testcase
type TestCasePriority uint16

const (
	// P0 critical priority test case
	P0 TestCasePriority = 0
	// P1 high priority test case
	P1 TestCasePriority = 1
	// P2 medium priority test case
	P2 TestCasePriority = 2
	// P3 low priority test case
	P3 TestCasePriority = 3
)

// TestCaseScope scope for test case
type TestCaseScope string

const (
	// ClusterScoped cluster test case scope
	ClusterScoped TestCaseScope = "Cluster"
	//NamespaceScoped test case scoped for a namespace
	NamespaceScoped TestCaseScope = "Namespaced"
)

// Options options for TestCase
type Options struct {
	// Name oof the test case
	Name string
	// Priority of the test case
	Priority TestCasePriority
	// Scope defines what kind of permissions this test case needs
	Scope TestCaseScope
}

func (o Options) defaultVals() Options {
	if o.Scope == TestCaseScope("") {
		o.Scope = NamespaceScoped
	}
	if o.Priority < P0 || o.Priority > P3 {
		o.Priority = P0
	}
	return o
}

// TestCaseBuilder builder for TestCases
// helps provide methods to construct
type TestCaseBuilder struct {
	opts     Options
	testFunc TestFunction
}

// TestCase constructor for test cases
func TestCase(opts Options) *TestCaseBuilder {
	return &TestCaseBuilder{
		opts: opts.defaultVals(),
	}
}

// P0Case builds a P0 case
func P0Case(name string) *TestCaseBuilder {
	return TestCase(Options{Name: name, Priority: P0})
}

// P1Case builds a P1 case
func P1Case(name string) *TestCaseBuilder {
	return TestCase(Options{Name: name, Priority: P1})
}

// P2Case builds a P1 case
func P2Case(name string) *TestCaseBuilder {
	return TestCase(Options{Name: name, Priority: P2})
}

// P3Case builds a P1 case
func P3Case(name string) *TestCaseBuilder {
	return TestCase(Options{Name: name, Priority: P3})
}

// WithFunc replaces the function with another given function
func (b *TestCaseBuilder) WithFunc(tc TestFunction) *TestCaseBuilder {
	b.testFunc = tc
	return b
}

// WithPriority sets priorities
func (b *TestCaseBuilder) WithPriority(prior TestCasePriority) *TestCaseBuilder {
	b.opts.Priority = prior
	b.opts = b.opts.defaultVals()
	return b
}

// Namespaced set the scope of the testcase as namespaced
func (b *TestCaseBuilder) Namespaced() *TestCaseBuilder {
	b.opts.Scope = NamespaceScoped
	return b
}

// Cluster set the scope of the testcase as a cluster scoped
func (b *TestCaseBuilder) Cluster() *TestCaseBuilder {
	b.opts.Scope = ClusterScoped
	return b
}

// P0 sets as P0
func (b *TestCaseBuilder) P0() *TestCaseBuilder {
	return b.WithPriority(P0)
}

// P1 sets as P1
func (b *TestCaseBuilder) P1() *TestCaseBuilder {
	return b.WithPriority(P1)
}

// P2 sets as P2
func (b *TestCaseBuilder) P2() *TestCaseBuilder {
	return b.WithPriority(P2)
}

// P3 sets as P3
func (b *TestCaseBuilder) P3() *TestCaseBuilder {
	return b.WithPriority(P3)
}

// Do builds and returns the test case
func (b *TestCaseBuilder) Do() bool {
	ctx := TestContext{
		Opts: b.opts,
	}
	fullName := fmt.Sprintf("[P%d][%s][%s]", b.opts.Priority, b.opts.Scope, b.opts.Name)
	return ginkgo.Describe(fullName, func() {
		ginkgo.By("Initializing " + fullName)
		ctx.Config = fmw.Config
		ctx.Context = fmw.Context
		ctx.SugaredLogger = fmw.SugaredLogger.Named(fullName)
		b.testFunc(ctx)
	})
}
