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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	// NamespaceScoped test case scoped for a Namespace
	NamespaceScoped TestCaseScope = "Namespaced"
)

// TestCaseLabel label for test case
type TestCaseLabel = string

const (
	ControllerLabel TestCaseLabel = "controller"
	WebhookLabel    TestCaseLabel = "webhook"
	WebServiceLabel TestCaseLabel = "webservice"
	CliLabel        TestCaseLabel = "cli"
)

// Options options for TestCase
type Options struct {
	// Name oof the test case
	Name string
	// Priority of the test case
	Priority TestCasePriority
	// Scope defines what kind of permissions this test case needs
	Scope TestCaseScope
	// Labels used to filter test cases when executing testing
	Labels []string
	// Condition used to check condition before testing
	Conditions []Condition
	// TestContextOptions used to setup TestContext
	TestContextOptions []TestContextOption
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

func (o Options) appendLabels(labels ...string) Options {
	m := make(map[string]struct{}, len(o.Labels))
	for _, label := range o.Labels {
		m[label] = struct{}{}
	}
	for _, newLabel := range labels {
		if _, exist := m[newLabel]; exist {
			continue
		}
		m[newLabel] = struct{}{}
		o.Labels = append(o.Labels, newLabel)
	}
	return o
}

func (o Options) checkCondition(testCtx *TestContext) error {
	for _, condition := range o.Conditions {
		if condition == nil {
			continue
		}
		if err := condition.Condition(testCtx); err != nil {
			return fmt.Errorf("condition %s check failed: %w", reflectName(condition), err)
		}
	}
	return nil
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

// WithLabels sets labels
func (b *TestCaseBuilder) WithLabels(labels ...string) *TestCaseBuilder {
	b.opts = b.opts.appendLabels(labels...)
	return b
}

// WithCondition sets conditions
func (b *TestCaseBuilder) WithCondition(funcs ...Condition) *TestCaseBuilder {
	b.opts.Conditions = append(b.opts.Conditions, funcs...)
	return b
}

// WithTestContextOptions sets options of TestContext
func (b *TestCaseBuilder) WithTestContextOptions(options ...TestContextOption) *TestCaseBuilder {
	b.opts.TestContextOptions = append(b.opts.TestContextOptions, options...)
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

// setupTestContext initial TestContext configuration
func (b *TestCaseBuilder) setupTestContext(ctx *TestContext) {
	ctx.Context = fmw.Context
	ctx.Config = fmw.Config
	ctx.SugaredLogger = fmw.SugaredLogger.Named(b.caseName())
	ctx.Scheme = fmw.Scheme

	for _, option := range b.opts.TestContextOptions {
		option(ctx)
	}

	if ctx.Namespace == "" {
		ctx.Namespace = "e2e-test-ns" + rand.String(5)
	}

	if ctx.Client == nil {
		c, err := client.New(ctx.Config, client.Options{Scheme: ctx.Scheme})
		Expect(err).To(Succeed())
		ctx.Client = c
	}
}

func (b *TestCaseBuilder) caseName() string {
	return fmt.Sprintf("[P%d][%s][%s]", b.opts.Priority, b.opts.Scope, b.opts.Name)
}

// Do build and return the test case
func (b *TestCaseBuilder) Do() bool {
	fullName := b.caseName()
	return Describe(fullName, Ordered, Labels(b.opts.Labels), func() {
		var testCtx = &TestContext{}

		BeforeAll(func() {
			b.setupTestContext(testCtx)

			// builtin condition
			conditions := append([]Condition{}, builtInConditions...)
			b.opts.Conditions = append(conditions, b.opts.Conditions...)
			err := b.opts.checkCondition(testCtx)
			if err != nil && skipWhenConditionMismatch == "true" {
				skipMsg := fmt.Sprintf("Skip test case, name: %s, reason: %s", fullName, err.Error())
				Skip(skipMsg)
			} else {
				Expect(err).To(Succeed())
			}
		})

		if b.testFunc != nil {
			b.testFunc(testCtx)
		}
	})
}

// DoFunc build and return the test case, just like the Do function
func (b *TestCaseBuilder) DoFunc(f TestFunction) bool {
	b.testFunc = f
	return b.Do()
}
