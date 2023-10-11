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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -package=multicluster -destination=../testing/mock/github.com/katanomi/pkg/multicluster/interface.go  github.com/katanomi/pkg/multicluster Interface

// Interface interface for a multi-cluster functionality
type Interface interface {
	GetConfig(ctx context.Context, clusterRef *corev1.ObjectReference) (config *rest.Config, err error)
	GetDynamic(ctx context.Context, clusterRef *corev1.ObjectReference) (dyn dynamic.Interface, err error)
	GetConfigFromCluster(ctx context.Context, cluster *unstructured.Unstructured) (config *rest.Config, err error)

	// ListClustersNamespaces lists all namespaces in all clusters
	// TODO: add this method to the interface and implementation
	ListClustersNamespaces(ctx context.Context, namespace string) (clusterNamespaces map[*corev1.ObjectReference][]corev1.Namespace, err error)
	// StartWarmUpClientCache used to start warming the client cache, only needs to be called once.
	StartWarmUpClientCache(ctx context.Context)

	// ClientGetter for getting client for a clusterRef and given scheme
	ClientGetter

	// NamespaceClustersGetter for getting list of clusters related by special namespace
	NamespaceClustersGetter
}

// NamespaceClustersGetter interface get list of clusters related by special namespace
type NamespaceClustersGetter interface {
	GetNamespaceClusters(ctx context.Context, namespace string) ([]corev1.ObjectReference, error)
}

// ClientGetter interface get client for a clusterRef and given scheme
type ClientGetter interface {
	GetClient(ctx context.Context, clusterRef *corev1.ObjectReference, scheme *runtime.Scheme) (clt client.Client, err error)
}
