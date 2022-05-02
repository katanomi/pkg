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

## TestCases

Most test cases can be written in a `ginkgo` fashion with a few helper methods to speedup construction and common logic. Cases can be started with:


   1. `TestCase`: constructor that receives an `Options` struct with all options for test case.
   1. `P0Case`: constructor method to set a name and a priority. Other levels are also available: `P1Case`, `P2Case`, `P3Case`.

After constructing a few more methods:

   1. `WithFunc`: takes a `TestFunction` that is given a context in which the test case executes.
   1. `Do`: finalizes the test case construction




```golang
package another

import (
	. "github.com/katanomi/pkg/testing/framework"
	. "github.com/onsi/ginkgo"
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


## TODO:

 - [ ]  Add support for traits/external dependencies: Add functions to declare dependencies of other systems/toolings and having independent implementations for each of them.

