/*
Copyright 2023 The AlaudaDevops Authors.

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
	"context"
	"fmt"

	"github.com/AlaudaDevops/pkg/testing"
	. "github.com/AlaudaDevops/pkg/testing/framework/base"
	. "github.com/onsi/ginkgo/v2"
)

func newCase(name string, priority TestCasePriority) *CaseBuilder {
	return &CaseBuilder{
		TestCaseBuilder{
			Name:     name,
			Priority: priority,
			TestContextGetter: TestContextGetterFunc(func() *TestContext {
				return &TestContext{
					Context:       fmw.Context,
					SugaredLogger: fmw.SugaredLogger.Named(name),
				}
			}),
			// set the unique case name label
			Labels: testing.Case(name),
		},
	}
}

// P0Case builds a P0 case
func P0Case(name string) *CaseBuilder {
	return newCase(name, P0)
}

// P1Case builds a P1 case
func P1Case(name string) *CaseBuilder {
	return newCase(name, P1)
}

// P2Case builds a P1 case
func P2Case(name string) *CaseBuilder {
	return newCase(name, P2)
}

// P3Case builds a P1 case
func P3Case(name string) *CaseBuilder {
	return newCase(name, P3)
}

// CaseBuilder builder for TestCases
type CaseBuilder struct {
	TestCaseBuilder
}

func (b *CaseBuilder) defaultVals() {
	if b.Priority < P0 || b.Priority > P3 {
		b.Priority = P0
	}
}

func (b *CaseBuilder) appendLabels(labels ...string) {
	m := make(map[string]struct{}, len(b.Labels))
	for _, label := range b.Labels {
		m[label] = struct{}{}
	}
	for _, newLabel := range labels {
		if _, exist := m[newLabel]; exist {
			continue
		}
		m[newLabel] = struct{}{}
		b.Labels = append(b.Labels, newLabel)
	}
}

// DoNotSkip set the test case to not skip when condition check failed
func (b *CaseBuilder) DoNotSkip() *CaseBuilder {
	b.FailedWhenConditionMismatch = false
	return b
}

// AllowSkip set the test case to skip when condition check failed
func (b *CaseBuilder) AllowSkip() *CaseBuilder {
	b.FailedWhenConditionMismatch = true
	return b
}

// WithFunc replaces the function with another given function
func (b *CaseBuilder) WithFunc(tc TestSpecFunc) *CaseBuilder {
	b.TestSpec = tc
	return b
}

// WithPriority sets priorities
func (b *CaseBuilder) WithPriority(prior TestCasePriority) *CaseBuilder {
	b.Priority = prior
	b.defaultVals()
	return b
}

// WithLabels sets labels
func (b *CaseBuilder) WithLabels(labels ...interface{}) *CaseBuilder {
	for _, label := range labels {
		switch label.(type) {
		case string:
			b.appendLabels(Labels{label.(string)}...)
		case Labels:
			b.appendLabels(label.(Labels)...)
		case interface{ Labels() Labels }:
			b.appendLabels(label.(interface{ Labels() Labels }).Labels()...)
		default:
			GinkgoWriter.Printf("unknown label type %s", label)
		}
	}
	return b
}

// WithCondition sets conditions
func (b *CaseBuilder) WithCondition(funcs ...Condition) *CaseBuilder {
	b.Conditions = append(b.Conditions, funcs...)
	return b
}

// P0 sets as P0
func (b *CaseBuilder) P0() *CaseBuilder {
	return b.WithPriority(P0)
}

// P1 sets as P1
func (b *CaseBuilder) P1() *CaseBuilder {
	return b.WithPriority(P1)
}

// P2 sets as P2
func (b *CaseBuilder) P2() *CaseBuilder {
	return b.WithPriority(P2)
}

// P3 sets as P3
func (b *CaseBuilder) P3() *CaseBuilder {
	return b.WithPriority(P3)
}

// Do build and return the test case
func (b *CaseBuilder) Do() bool {
	return b.DoWithContext(context.Background())
}

// DoWithContext run a test case with special context
// The context can be used for construct case layouts
func (b *CaseBuilder) DoWithContext(ctx context.Context) bool {
	fullName := b.CaseName()
	return Describe(fullName, Ordered, Labels(b.Labels), func() {
		var testCtx = &TestContext{Context: ctx}

		BeforeAll(func() {
			*testCtx = *b.GetTestContext()
			skip, err := b.CheckCondition(testCtx)
			if err != nil {
				if skip {
					Skip(fmt.Sprintf("Skip test case, reason: %s", err.Error()))
				} else {
					Fail(fmt.Sprintf("Test case failed, reason %s", err.Error()))
				}
			}
		})

		if b.TestSpec != nil {
			b.TestSpec(testCtx)
		}
	})
}

// DoFunc build and return the test case, just like the Do function
func (b *CaseBuilder) DoFunc(f TestSpecFunc) bool {
	b.TestSpec = f
	return b.Do()
}
