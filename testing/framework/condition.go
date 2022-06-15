/*
Copyright 2022 The Katanomi Authors.

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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	skipWhenConditionMismatch = GetDefaultEnv("SKIP_WHEN_CONDITION_MISMATCH", "true")
)

var builtInConditions = []Condition{
	&TestNamespaceCondition{},
}

// ConditionFunc helper function to wrapping condition
type ConditionFunc func(testCtx *TestContext) error

// Condition implement the Condition interface
func (c ConditionFunc) Condition(testCtx *TestContext) error {
	return c(testCtx)
}

// Condition describe the conditions which the test must match
type Condition interface {
	Condition(testCtx *TestContext) error
}

// TestNamespaceCondition generate namespace for testing
type TestNamespaceCondition struct{}

// Condition implement the Condition interface
// Delete the namespace when it already exists, then create a new one.
// After the testing is completed, delete the namespace as well.
func (t *TestNamespaceCondition) Condition(testCtx *TestContext) error {
	var (
		clt = testCtx.Client
		ctx = testCtx.Context
		err error
	)

	ns := v1.Namespace{}
	key := types.NamespacedName{Name: testCtx.Namespace}
	err = clt.Get(ctx, key, &ns)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	if ns.Name != "" {
		if err = clt.Delete(ctx, &ns); err != nil {
			return err
		}
		MustRollback(testCtx, &ns)
	}

	ns = v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: testCtx.Namespace},
	}
	if err = clt.Create(ctx, &ns); err != nil {
		return err
	}

	DeferCleanup(func() error {
		MustRollback(testCtx, &ns)
		return nil
	})
	return nil
}

// NewReconcileCondition helper function for constructing a `ReconcileCondition` object
func NewReconcileCondition(obj client.Object, objCheckFun func(obj client.Object) bool) *ReconcileCondition {
	c := &ReconcileCondition{
		Poller:      &Poller{},
		obj:         obj,
		objCheckFun: objCheckFun,
	}

	return c
}

// ReconcileCondition check if a controller is available
type ReconcileCondition struct {
	*Poller
	obj         client.Object
	objCheckFun func(obj client.Object) bool
}

// WithPoller customize poller settings
func (n *ReconcileCondition) WithPoller(interval, timeout time.Duration) *ReconcileCondition {
	n.Interval = interval
	n.Timeout = timeout

	return n
}

// Condition implement the Condition interface
// Apply an object and wait for reconciliation, then check the status via `objCheckFun`
func (n *ReconcileCondition) Condition(testCtx *TestContext) error {
	var (
		clt = testCtx.Client
		ctx = testCtx.Context
		err error
	)

	n.obj.SetNamespace(testCtx.Namespace)

	By(fmt.Sprintf("try to create %s resource", n.obj.GetName()))
	err = clt.Create(ctx, n.obj)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	if err == nil {
		defer func() {
			MustRollback(testCtx, n.obj)
		}()
	}

	By(fmt.Sprintf("wait for %s to be ready", n.obj.GetName()))
	brokerKey := types.NamespacedName{Namespace: n.obj.GetNamespace(), Name: n.obj.GetName()}
	interval, timeout := n.Settings()
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		if err = clt.Get(ctx, brokerKey, n.obj); err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			} else {
				return false, err
			}
		}
		if n.objCheckFun != nil && !n.objCheckFun(n.obj) {
			return false, nil
		}
		return true, nil
	})
}

// MustRollback delete an object and wait for completion
func MustRollback(testCtx *TestContext, obj client.Object, opts ...client.DeleteOption) {
	Expect(testCtx.Client).NotTo(BeNil())
	Expect(testCtx.Context).NotTo(BeNil())
	Expect(obj).NotTo(BeNil())

	err := testCtx.Client.Delete(testCtx.Context, obj, opts...)
	Expect(err).To(Succeed())

	WaitRollback(testCtx, obj)
}

// WaitRollback Wait for the delete object behavior to complete
func WaitRollback(testCtx *TestContext, obj client.Object) {
	key := types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()}
	Eventually(func(g Gomega) error {
		err := testCtx.Client.Get(testCtx.Context, key, obj)
		g.Expect(errors.IsNotFound(err)).To(BeTrue())
		return nil
	}).WithPolling(time.Second).WithTimeout(time.Minute).Should(Succeed())
}
