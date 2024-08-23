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

	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	installed := InstallCheck(ctx, clt, &pipev1beta1.TaskList{}, &client.ListOptions{})
	if !installed {
		return false
	}
	// If the PipelineRun does not exist, it is likely that the TaskRun does not exist either.
	return InstallCheck(ctx, clt, &pipev1beta1.PipelineRunList{}, &client.ListOptions{})
}

// InstallCheck common check component installed method.
func InstallCheck(ctx context.Context, clt client.Client, value client.ObjectList, opts ...client.ListOption) bool {
	log := logging.FromContext(ctx)
	if err := clt.List(ctx, value, opts...); err != nil {
		log.Debugw("list operation failed", "err", err)
		return false
	}
	return true
}
