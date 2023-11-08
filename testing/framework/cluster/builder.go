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

package cluster

import (
	"fmt"

	"github.com/katanomi/pkg/multicluster"
	"github.com/katanomi/pkg/testing"
	"github.com/katanomi/pkg/testing/framework/base"
	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/runtime"
	pkgrand "k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/pkg/injection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ControllerLabel base.TestCaseLabel = "controller"
	WebhookLabel    base.TestCaseLabel = "webhook"
	WebServiceLabel base.TestCaseLabel = "webservice"
	CliLabel        base.TestCaseLabel = "cli"
)

// TestCaseScope scope for test case
type TestCaseScope string

const (
	// ClusterScoped cluster test case scope
	ClusterScoped TestCaseScope = "Cluster"
	// NamespaceScoped test case scoped for a Namespace
	NamespaceScoped TestCaseScope = "Namespaced"
)

func NewTestCaseBuilder(baseBuilder base.TestCaseBuilder) *TestCaseBuilder {
	return &TestCaseBuilder{
		baseBuilder: baseBuilder,
	}
}

type TestCaseBuilder struct {
	baseBuilder base.TestCaseBuilder

	// Conditions condition list which will be checked before test case execution
	Conditions []Condition

	// Scope defines what kind of permissions this test case needs
	Scope TestCaseScope

	// TestContextOptions used to setup TestContext
	TestContextOptions []TestContextOption

	// Scheme the scheme for initializing the k8s client
	Scheme *runtime.Scheme

	testSpecFunc TestSpecFunc
}

// WithCondition sets conditions
func (b *TestCaseBuilder) WithCondition(funcs ...Condition) *TestCaseBuilder {
	b.Conditions = append(b.Conditions, funcs...)
	return b
}

// WithScheme adds a scheme object to the framework
func (b *TestCaseBuilder) WithScheme(scheme *runtime.Scheme) *TestCaseBuilder {
	b.Scheme = scheme
	return b
}

// WithTestContextOptions sets options of TestContext
func (b *TestCaseBuilder) WithTestContextOptions(options ...TestContextOption) *TestCaseBuilder {
	b.TestContextOptions = append(b.TestContextOptions, options...)
	return b
}

// WithNamespacedScope set the scope of the testcase as namespaced
func (b *TestCaseBuilder) WithNamespacedScope() *TestCaseBuilder {
	b.Scope = NamespaceScoped
	return b
}

// WithClusterScope set the scope of the testcase as a cluster scoped
func (b *TestCaseBuilder) WithClusterScope() *TestCaseBuilder {
	b.Scope = ClusterScoped
	return b
}

// setupTestCase initial TestContext configuration
func (b *TestCaseBuilder) setupTestContext() *TestContext {
	ctx := &TestContext{
		TestContext: *b.baseBuilder.GetTestContext(),
	}

	cfg := ctrl.GetConfigOrDie()
	ctx.Config = cfg
	ctx.Context = injection.WithConfig(ctx.Context, cfg)

	if b.Scheme == nil {
		b.Scheme = FromSharedScheme(ctx.Context)
	}
	ctx.Scheme = b.Scheme

	for _, option := range b.TestContextOptions {
		option(ctx)
	}

	if ctx.Namespace == "" {
		ctx.Namespace = "e2e-test-ns" + pkgrand.String(5)
	}

	if ctx.Client == nil {
		c, err := client.New(ctx.Config, client.Options{Scheme: ctx.Scheme})
		if err != nil {
			panic(err)
		}
		ctx.Client = c
	}

	if ctx.MultiClusterClient == nil {
		mc, err := multicluster.NewClusterRegistryClient(ctx.Config)
		if err != nil {
			panic(err)
		}
		ctx.MultiClusterClient = mc
	}

	return ctx
}

func (b *TestCaseBuilder) checkCondition(testContext *TestContext) (skip bool, err error) {
	skip, err = b.baseBuilder.CheckCondition(&testContext.TestContext)
	if err != nil {
		return
	}

	conditions := append([]Condition{&TestNamespaceCondition{}}, b.Conditions...)
	for _, condition := range conditions {
		if condition == nil {
			continue
		}
		if err = condition.Condition(testContext); err != nil {
			err = fmt.Errorf("condition %s check failed: %w", testing.ReflectName(condition), err)
			break
		}
	}
	if err != nil && !b.baseBuilder.FailedWhenConditionMismatch {
		skip = true
	}

	return
}

// WithFunc replaces the function with another given function
func (b *TestCaseBuilder) WithFunc(tc TestSpecFunc) *TestCaseBuilder {
	b.testSpecFunc = tc
	return b
}

// Do build and return the test case
func (b *TestCaseBuilder) Do() bool {
	fullName := b.baseBuilder.CaseName()
	return Describe(fullName, Ordered, Labels(b.baseBuilder.Labels), func() {
		var testCtx = &TestContext{}

		BeforeAll(func() {
			*testCtx = *b.setupTestContext()
			skip, err := b.checkCondition(testCtx)
			if err != nil {
				if skip {
					Skip(fmt.Sprintf("Skip test case, reason: %s", err.Error()))
				} else {
					Fail(fmt.Sprintf("Test case failed, reason %s", err.Error()))
				}
			}
		})

		if b.testSpecFunc != nil {
			b.testSpecFunc(testCtx)
		}
	})
}

// DoFunc build and return the test case, just like the Do function
func (b *TestCaseBuilder) DoFunc(f TestSpecFunc) bool {
	b.testSpecFunc = f
	return b.Do()
}
