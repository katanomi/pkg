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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrNilReference = errors.New("nil reference for clusterRef object")
var ErrNoNamespaceProvided = errors.New("namespace must be provided")
var ErrNoNameProvided = errors.New("name must be provided")
var ErrDoesNotHaveEndpoints = errors.New("cluster object does not have spec.kubernetesApiEndpoints.serverEndpoints")
var ErrDoesNotHaveServerAddress = errors.New("cluster object does not have spec.kubernetesApiEndpoints.serverEndpoints.serverAddress")
var ErrDoesNotHaveToken = errors.New("secret does not have data.token")

// ClusterRegistryClient implements the deprecated cluster registry cluster resource multi cluster client
// https://github.com/kubernetes-retired/cluster-registry/blob/master/pkg/apis/clusterregistry/v1alpha1/types.go
type ClusterRegistryClient struct {
	dynamic.Interface
}

var _ Interface = &ClusterRegistryClient{}

// NewClusterRegistryClient initiates a ClusterRegistryClient
func NewClusterRegistryClient(config *rest.Config) (Interface, error) {
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ClusterRegistryClient{Interface: dyn}, nil
}

// NewClusterRegistryClientOrDie initiates a ClusterRegistryClient and
// panics if it fails
func NewClusterRegistryClientOrDie(config *rest.Config) Interface {
	clt, err := NewClusterRegistryClient(config)
	if err != nil {
		panic(err)
	}
	return clt
}

var ClusterRegistryGroupVersion = schema.GroupVersion{Group: "clusterregistry.k8s.io", Version: "v1alpha1"}
var ClusterRegistryGVK = ClusterRegistryGroupVersion.WithKind("Cluster")

// GetConfig returns the configuration based on the Cluster
func (m *ClusterRegistryClient) GetConfig(ctx context.Context, clusterRef *corev1.ObjectReference) (config *rest.Config, err error) {
	if err = m.validateRef(clusterRef); err != nil {
		return
	}
	var cluster *unstructured.Unstructured
	if cluster, err = m.getClusterByRef(ctx, clusterRef); err != nil {
		return
	}
	config, err = m.getConfigFromCluster(ctx, cluster)
	return
}

// GetClient returns a client using the cluster configuration
func (m *ClusterRegistryClient) GetClient(ctx context.Context, clusterRef *corev1.ObjectReference, scheme *runtime.Scheme) (clt client.Client, err error) {
	config, configErr := m.GetConfig(ctx, clusterRef)
	if configErr != nil {
		err = configErr
		return
	}
	clt, err = client.New(config, client.Options{Scheme: scheme})
	return
}

// GetDynamic returns a dynamic client using the cluster configuration
func (m *ClusterRegistryClient) GetDynamic(ctx context.Context, clusterRef *corev1.ObjectReference) (dyn dynamic.Interface, err error) {
	config, configErr := m.GetConfig(ctx, clusterRef)
	if configErr != nil {
		err = configErr
		return
	}
	dyn, err = dynamic.NewForConfig(dynamic.ConfigFor(config))
	return
}

// TODO: Change to use kubernets error types
func (m *ClusterRegistryClient) validateRef(clusterRef *corev1.ObjectReference) (err error) {
	if clusterRef == nil {
		return ErrNilReference
	}
	if clusterRef.Name == "" {
		return ErrNoNameProvided
	}
	if clusterRef.Namespace == "" {
		return ErrNoNamespaceProvided
	}
	return nil
}

func (m *ClusterRegistryClient) getClusterByRef(ctx context.Context, clusterRef *corev1.ObjectReference) (cluster *unstructured.Unstructured, err error) {
	// this method aims to get a specific cluster registry implementation
	cluster, err = m.Interface.
		Resource(ClusterRegistryGroupVersion.WithResource("clusters")).
		Namespace(clusterRef.Namespace).
		Get(ctx, clusterRef.Name, metav1.GetOptions{})
	return
}

func (m *ClusterRegistryClient) getConfigFromCluster(ctx context.Context, cluster *unstructured.Unstructured) (config *rest.Config, err error) {
	if cluster == nil {
		err = ErrNilReference
		return
	}

	var jsonData []byte
	jsonData, err = json.Marshal(cluster)
	if err != nil {
		return
	}

	clusterObj := &Cluster{}
	if err = json.Unmarshal(jsonData, clusterObj); err != nil {
		return
	}

	if len(clusterObj.Spec.KubernetesAPIEndpoints.ServerEndpoints) == 0 {
		err = ErrDoesNotHaveEndpoints
		return
	}
	address := clusterObj.Spec.KubernetesAPIEndpoints.ServerEndpoints[0].ServerAddress
	caBundle := clusterObj.Spec.KubernetesAPIEndpoints.CABundle

	config = &rest.Config{
		Host: address,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caBundle,
		},
	}

	// Uses the controller secret for controller auth
	// this configuration or config should only be used by controllers
	if clusterObj.Spec.AuthInfo.Controller != nil {
		name := clusterObj.Spec.AuthInfo.Controller.Name
		namespace := clusterObj.Spec.AuthInfo.Controller.Namespace
		if namespace == "" {
			namespace = cluster.GetNamespace()
		}
		secretObj, secretErr := m.Interface.Resource(corev1.SchemeGroupVersion.WithResource("secrets")).
			Namespace(namespace).
			Get(ctx, name, metav1.GetOptions{})
		if secretErr != nil {
			err = secretErr
			return
		}
		token, ok, secretErr := unstructured.NestedString(secretObj.Object, "data", "token")
		fmt.Println("data.token", token, "ok", ok)
		if secretErr != nil {
			err = secretErr
			return
		} else if !ok {
			err = ErrDoesNotHaveToken
			return
		}
		tokenBytes, secretErr := base64.StdEncoding.DecodeString(token)
		if secretErr != nil {
			err = secretErr
			return
		}
		config.BearerToken = string(tokenBytes)
	}
	return
}
