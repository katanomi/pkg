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
	"context"
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

	// ObjectLists the list of objects to be reported in the test case
	// whenever a test case fails
	ObjectLists []client.ObjectList

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

// WithObjectListReport will add a list of objects to be reported just after each test case
// if the test case fails.
// Can be used with testing.ListByGVK() to quickly create an ObjectList.
// Example:
// var _ = P0Case("test case").
//
//	Cluster().
//	WithObjectListReport(testing.ListByGVK(schema.GroupVersionKind{
//	  Group:   "tekton.dev",
//	  Version: "v1beta1",
//	  Kind:    "PipelineRunList",
//	}, &corev1.ConfigMapList{}).DoWithContext(ctx)
//
// Note that multiple ObjectList can be added to the same test case.
func (b *TestCaseBuilder) WithObjectListReport(objectList ...client.ObjectList) *TestCaseBuilder {
	b.ObjectLists = append(b.ObjectLists, objectList...)
	return b
}

// WithFunc replaces the function with another given function
func (b *TestCaseBuilder) WithFunc(tc TestSpecFunc) *TestCaseBuilder {
	b.testSpecFunc = tc
	return b
}

// DoWithContext runs the test case builder with the given context.
// It sets up the test context, checks test conditions, adds object reporting,
// and runs the test spec function if provided.
func (b *TestCaseBuilder) DoWithContext(ctx context.Context) bool {
	fullName := b.baseBuilder.CaseName()
	return Describe(fullName, Ordered, Labels(b.baseBuilder.Labels), func() {
		var testCtx = &TestContext{
			TestContext: base.TestContext{
				Context: ctx,
			},
		}

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

		if len(b.ObjectLists) > 0 {
			JustAfterEach(func() {
				ReportObjectsLists(testCtx, b.ObjectLists...)
			})
		}

		if b.testSpecFunc != nil {
			b.testSpecFunc(testCtx)
		}
	})
}

// Do build and return the test case
func (b *TestCaseBuilder) Do() bool {
	return b.DoWithContext(context.Background())
}

// DoFunc build and return the test case, just like the Do function
func (b *TestCaseBuilder) DoFunc(f TestSpecFunc) bool {
	b.testSpecFunc = f
	return b.Do()
}
