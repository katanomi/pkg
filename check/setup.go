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

const (
	// NotExistNamespace returns not exist namespace name.
	NotExistNamespace = "not exist namespace"
)

// InstalledTekton Check if tekton is installed
func InstalledTekton(ctx context.Context, clt client.Client) bool {
	return InstallCheck(ctx, clt, &pipev1beta1.PipelineRunList{})
}

// InstallCheck common check component installed method.
func InstallCheck(ctx context.Context, clt client.Client, value client.ObjectList) bool {
	log := logging.FromContext(ctx)
	if err := clt.List(ctx, value, client.InNamespace(NotExistNamespace)); err != nil {
		log.Debugw("list operation failed", "err", err)
		return false
	}
	return true
}
