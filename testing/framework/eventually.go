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
	"github.com/katanomi/pkg/testing"
	. "github.com/onsi/gomega"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsReadyCheckFunc additional function to check if object is ready in an Eventually context
/* example:
   ```
   func(g Gomega, obj client.Object) error {
	  if obj.GetName() == "" {
		fmt.Errorf("should have a name")
	  }
      // or use g.Expect instead
      g.Expect(obj.GetName()).To(Equal(""))
	  return nil
   }
   ```
*/
type IsReadyCheckFunc func(g Gomega, obj client.Object) error

type hasGetTopLevelCondition interface {
	GetTopLevelCondition() *apis.Condition
}

// ResourceTopConditionIsReadyEventually generic function to check if a resource is ready
// resource type must implement  GetTopLevelCondition() *apis.Condition function
func ResourceTopConditionIsReadyEventually(ctx *TestContext, obj client.Object, readyFuncs ...IsReadyCheckFunc) func(g Gomega) error {
	key := client.ObjectKeyFromObject(obj)
	return func(g Gomega) error {
		condManager, ok := obj.(hasGetTopLevelCondition)
		g.Expect(ok).To(BeTrue(), "object should implement GetTopLevelCondition() *apis.Condition function")

		g.Expect(ctx.Client.Get(ctx.Context, key, obj)).To(Succeed(), "should get the object")

		if err := testing.ConditionIsReady(condManager.GetTopLevelCondition()); err != nil {
			return err
		}
		for _, isReady := range readyFuncs {
			if err := isReady(g, obj); err != nil {
				return err
			}
		}
		return nil
	}
}
