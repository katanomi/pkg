# e2e framework

Enables to easily build and start a number of automated test cases. Depends on ginkgo to describe test suites.

A sample usage is found in [examples](../../examples/sample-e2e)

## Framework

Holds the global configurations for the framework. It mainly:

 1. `*rest.Config` Kubernetes client configuration
 1. `context.Context` a common context
 1. `*zap.SugaredLogger` a logger for test cases

 etc.


```golang


import (
	"testing"

    // Adds test cases from other packages here
	_ "github.com/katanomi/pkg/examples/sample-e2e/another"
	"github.com/katanomi/pkg/testing/framework"
)

var fmw = framework.New("sample-e2e")

func TestMain(m *testing.M) {
	fmw.SynchronizedBeforeSuite(nil).
		SynchronizedAfterSuite(nil).
		MRun(m)
}

func TestE2E(t *testing.T) {
	// start step to run e2e
	fmw.Run(t)
}

```

## Configuration

Each project e2e test may depend on some configurations. By default, if these configurations are missing, the test will not execute normally.

The framework has a built-in NewConfigCondition tool method to help us load and manage configurations. By default configuration files are stored in the `katanomi-e2e` namespace of the target cluster, and the name of the configmap will be prefixed with `e2e-config`.

```go
func NewConfigCondition(configName string, obj interface{}) *configCondition {
	c := &configCondition{
		name: configName,
		obj:  obj,
	}

	return c
}
```

Of course, you can use environment variables to change these defaults:
```go
E2E_CONFIG_NAMESPACE
E2E_CONFIG_NAME_PREFIX
```

## TestCases

Most test cases can be written in a `ginkgo` fashion with a few helper methods to speedup construction and common logic. Cases can be started with:


   1. `TestCase`: constructor that receives an `Options` struct with all options for test case.
   1. `P0Case`: constructor method to set a name and a priority. Other levels are also available: `P1Case`, `P2Case`, `P3Case`.

After constructing a few more methods:

   1. `WithFunc`: takes a `TestFunction` that is given a context in which the test case executes.
   1. `Do`: finilizes the test case construction




```golang
package another

import (
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)



var _ = TestCase(Options{Name: "Testinge2e", Priority: P0, Scope: NamespaceScoped}).WithFunc(func(ctx TestContext) {
	BeforeEach(func() {
		ctx.Debugw("some debug message")
		// fmt.Println("TestCase BeforeEach", ctx.Config)
	})
	It("should succeed", func() {
		Expect(ctx.Config).ToNot(BeNil())
	})
}).Do()

var _ = P1Case("another-test").Cluster().WithFunc(func(ctx TestContext) {
	// test case
	BeforeEach(func() {
		ctx.Debugw("before each in another pkg")
	})
	AfterEach(func() {
		ctx.Debugw("after each in another pkg")
	})
	Context("With a cluster scoped test case", func() {
		JustBeforeEach(func() {
			ctx.Infow("just before each in another pkg")
		})
		JustAfterEach(func() {
			ctx.Infow("just after each in another pkg")
		})
		It("it", func() {
			Expect(ctx.Config).ToNot(BeNil())
		})
	})
}).Do()

```

