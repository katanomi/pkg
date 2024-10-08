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

// Package check contains functions to check component installed.
package check

import (
	"context"
	"time"

	pkgclient "github.com/katanomi/pkg/client"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// InstalledTekton Check if tekton is installed
func InstalledTekton(ctx context.Context, clt client.Client) bool {
	// Checks whether the task resource exists in the cluster.

	// Since the new version of Tekton defaults to storing tasks in v1 version,
	// fetching tasks in v1beta1 version requires a conversion.
	// If the conversion fails, it also indicates that Tekton does not exist.

	// In the new version, when katanomi is deployed, the Task resource will be automatically imported, so the resource should exist.
	// Even if the resource does not exist, it's fine.
	// Although the list operation won't fail, an error will still occur when importing resources.
	// It also won't trigger the issue where client-go's cache fails to list the Task resource.
	installed := InstallationCheck(ctx, clt, &pipev1beta1.TaskList{}, &client.ListOptions{})
	if !installed {
		return false
	}

	// If the PipelineRun does not exist, it is likely that the TaskRun does not exist either.
	return InstallationCheck(ctx, clt, &pipev1beta1.PipelineRunList{}, &client.ListOptions{})
}

// InstallationCheck common check component installed method.
func InstallationCheck(ctx context.Context, clt client.Client, value client.ObjectList, opts ...client.ListOption) bool {
	log := logging.FromContext(ctx)

	defer func(startTime time.Time) {
		log.Debugw("Installation check", "gvk", value.GetObjectKind().GroupVersionKind(), "UsedTime", time.Since(startTime))
	}(time.Now())

	// We only need to check the first one to determine whether the crd exists and the component is ready.
	// If set `resourceVersion` to 0, it will ignore the `limit` option.
	opts = append(opts, pkgclient.NewListOptions().WithLimit(1).WithUnsafeDisableDeepCopy().Build())

	if err := clt.List(ctx, value, opts...); err != nil {
		log.Infow("list operation failed", "err", err)
		return false
	}

	if prList, ok := value.(*pipev1beta1.PipelineRunList); ok && prList != nil {
		log.Debugw("list PipelineRun success", "count", len(prList.Items))
	}

	// Set the group version kind to the object, otherwise the log will not show the GVK.
	gvk, err := apiutil.GVKForObject(value, clt.Scheme())
	if err != nil {
		log.Errorw("get GVK failed", "err", err)
	}
	value.GetObjectKind().SetGroupVersionKind(gvk)

	return true
}
