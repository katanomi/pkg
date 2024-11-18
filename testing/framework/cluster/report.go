/*
Copyright 2024 The Katanomi Authors.

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

	"github.com/AlaudaDevops/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ReportObjectsLists will dump objects into Current Ginkgo Report when test case failed.
// Example:
//
//	JustAfteEach(func() {
//	  ReportObjectsLists(ctx, &corev1.ConfigMapList{}, &corev1.PodList{})
//	})
func ReportObjectsLists(ctx *TestContext, objects ...ctrlclient.ObjectList) {
	// Adding InNamespace option will make test reports clear
	// and will avoid extra noise from unrelated namespaces.
	// For Cluster scoped resources this will have no effect and
	// will return all resources in the cluster.
	options := []ctrlclient.ListOption{}
	if ctx.Namespace != "" {
		options = append(options, ctrlclient.InNamespace(ctx.Namespace))
	}
	for _, _item := range objects {
		item := _item
		ReportObjectListWithOptions(ctx, item, options...)
	}
}

// ReportObjectListWithOptions dumps the provided object list into the
// current Ginkgo spec report if the test case has failed. It allows
// passing client list options.
// Example:
//
//	JustAfteEach(func() {
//	  ReportObjectListWithOptions(ctx, &corev1.ConfigMapList{}, client.InNamespace(ctx.Namespace))
//	})
//
// When using testing.ListByGVK:
//
//	JustAfteEach(func() {
//	  ReportObjectListWithOptions(ctx, ListByGVK(schema.GroupVersionKind{
//	    Group:   "katanomi.dev",
//	    Version: "v1alpha1",
//	    Kind:    "ClusterRuleList",
//	  }), client.InNamespace(ctx.Namespace))
//	})
func ReportObjectListWithOptions(ctx *TestContext, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) {
	name := fmt.Sprintf("%s", list.GetObjectKind().GroupVersionKind().String())
	err := ctx.Client.List(ctx.Context, list, opts...)
	if err != nil {
		GinkgoWriter.Printf("list error: %s, list: %#v", err.Error(), list)
		return
	}
	testing.AddReportEntryAsYaml(name, list)
}
