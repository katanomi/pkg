/*
Copyright 2021 The Katanomi Authors.

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

package multicluster

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Interface interface for a multi-cluster functionality
type Interface interface {
	GetConfig(ctx context.Context, clusterRef *corev1.ObjectReference) (config *rest.Config, err error)
	GetClient(ctx context.Context, clusterRef *corev1.ObjectReference, scheme *runtime.Scheme) (clt client.Client, err error)
	GetDynamic(ctx context.Context, clusterRef *corev1.ObjectReference) (dyn dynamic.Interface, err error)

	// TODO: add this method to the interface and implementation
	// ListClustersNamespaces(ctx context.Context, namespace string) (clusterNamespaces map[*corev1.ObjectReference][]corev1.Namespace , err error)
}
