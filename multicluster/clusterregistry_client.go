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
	"sync"

	"github.com/katanomi/pkg/parallel"
	"knative.dev/pkg/logging"

	"github.com/katanomi/pkg/tracing"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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

	insecure bool
	// proxy host for accessing cluster
	clusterProxyHost string
	// proxy host for accessing cluster, support {name} placeholder with the actual cluster name
	clusterProxyPath string
}

var _ Interface = &ClusterRegistryClient{}

// NewClusterRegistryClient initiates a ClusterRegistryClient
func NewClusterRegistryClient(config *rest.Config, options ...ClusterRegistryClientOption) (Interface, error) {
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	registryClient := &ClusterRegistryClient{Interface: dyn}
	for _, option := range options {
		option(registryClient)
	}

	return registryClient, nil
}

// NewClusterRegistryClientOrDie initiates a ClusterRegistryClient and
// panics if it fails
func NewClusterRegistryClientOrDie(config *rest.Config, options ...ClusterRegistryClientOption) Interface {
	clt, err := NewClusterRegistryClient(config, options...)
	if err != nil {
		panic(err)
	}
	return clt
}

// ClusterRegistryClientOption functions for configuring a ClusterRegistryClient
type ClusterRegistryClientOption func(*ClusterRegistryClient)

// ClusterProxyOption sets the proxy host and path for the cluster registry client
func ClusterProxyOption(proxyHost string, proxyPath string) ClusterRegistryClientOption {
	return func(c *ClusterRegistryClient) {
		c.clusterProxyHost = proxyHost
		c.clusterProxyPath = proxyPath
	}
}

// ClusterProxyInsecure allows specifying whether the client should use an insecure connection.
func ClusterProxyInsecure(insecure bool) ClusterRegistryClientOption {
	return func(c *ClusterRegistryClient) {
		c.insecure = insecure
	}
}

var ClusterRegistryGroupVersion = schema.GroupVersion{Group: "clusterregistry.k8s.io", Version: "v1alpha1"}
var ClusterRegistryGVK = ClusterRegistryGroupVersion.WithKind("Cluster")
var ClusterGVR = ClusterRegistryGroupVersion.WithResource("clusters")

// GetConfig returns the configuration based on the Cluster
func (m *ClusterRegistryClient) GetConfig(ctx context.Context, clusterRef *corev1.ObjectReference) (config *rest.Config, err error) {
	if err = m.validateRef(clusterRef); err != nil {
		return
	}
	var cluster *unstructured.Unstructured
	if cluster, err = m.getClusterByRef(ctx, clusterRef); err != nil {
		return
	}
	config, err = m.GetConfigFromCluster(ctx, cluster)
	if m.clusterProxyHost != "" {
		proxyHost, err := ClusterProxyHost(m.clusterProxyHost, m.clusterProxyPath, cluster.GetName())
		if err != nil {
			return nil, err
		}
		config.Host = proxyHost
	}

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

func (m *ClusterRegistryClient) GetConfigFromCluster(ctx context.Context, cluster *unstructured.Unstructured) (config *rest.Config, err error) {
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

	tlsConfig := rest.TLSClientConfig{
		Insecure: m.insecure,
	}
	// If connections are not intended to be insecure, configure CAData with the certificate bundle.
	// This ensures SSL certificates are verified unless explicitly disabled by setting the insecure flag.
	if !m.insecure {
		tlsConfig.CAData = caBundle
	}
	config = &rest.Config{
		Host:            address,
		TLSClientConfig: tlsConfig,
	}

	config.Wrap(tracing.WrapTransport)

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

// ListClustersNamespaces will list namespace with name "namespace" in all clusters
func (m *ClusterRegistryClient) ListClustersNamespaces(ctx context.Context, namespace string) (clusterNamespaces map[*corev1.ObjectReference][]corev1.Namespace, err error) {
	clusterRefs, err := m.GetNamespaceClusters(ctx, namespace)
	if err != nil {
		return nil, err
	}

	maxConcurrency := 10
	log := logging.FromContext(ctx)
	p := parallel.P(log, "ListClusterNamespace").SetConcurrent(maxConcurrency)

	resultMap := sync.Map{}

	for _, clusterRef := range clusterRefs {
		_clusterRef := clusterRef
		p.Add(func() (interface{}, error) {

			clusterClient, err := m.GetClient(ctx, &_clusterRef, clientgoscheme.Scheme)
			if err != nil {
				log.Errorw("error to get cluster client", "cluster", clusterRef.Namespace+"/"+clusterRef.Name, "err", err.Error())
				return nil, err
			}

			ns := &corev1.Namespace{}
			err = clusterClient.Get(ctx, client.ObjectKey{Name: namespace}, ns)
			if err != nil {
				return nil, err
			}

			resultMap.Store(clusterRef, ns)
			return nil, nil
		})
	}

	_, err = p.Do().Wait()
	if err != nil {
		return nil, err
	}

	clusterNamespaces = map[*corev1.ObjectReference][]corev1.Namespace{}
	resultMap.Range(func(key, value interface{}) bool {
		clusterRef, ok := key.(corev1.ObjectReference)
		if !ok {
			return true
		}
		ns, ok := value.(*corev1.Namespace)
		if !ok {
			return true
		}

		clusterNamespaces[&clusterRef] = []corev1.Namespace{*ns}
		return true
	})

	return clusterNamespaces, nil
}

// StartWarmUpClientCache used to start warming the client cache, only needs to be called once.
func (m *ClusterRegistryClient) StartWarmUpClientCache(ctx context.Context) {
}

// GetNamespaceClusters returns a list of clusters related by namespace
func (m *ClusterRegistryClient) GetNamespaceClusters(ctx context.Context, namespace string) (clusterRefs []corev1.ObjectReference, err error) {
	clusters, err := m.Interface.
		Resource(ClusterRegistryGroupVersion.WithResource("clusters")).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{ResourceVersion: "0"})
	if err != nil {
		return nil, err
	}
	clusterRefs = make([]corev1.ObjectReference, 0, len(clusters.Items))
	for _, cluster := range clusters.Items {
		clusterRef := corev1.ObjectReference{
			Kind:       cluster.GetKind(),
			Namespace:  cluster.GetNamespace(),
			Name:       cluster.GetName(),
			APIVersion: cluster.GetAPIVersion(),
		}
		clusterRefs = append(clusterRefs, clusterRef)
	}
	return clusterRefs, nil
}
