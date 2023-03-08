/*
Copyright 2023 The Katanomi Authors.

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

	"github.com/emicklei/go-restful/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	apiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/system"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewCrossClusterSubjectReview constructs a new CrossClusterSubjectReview
func NewCrossClusterSubjectReview(mClient Interface, scheme *runtime.Scheme, restMapper meta.RESTMapper) *CrossClusterSubjectReview {
	return &CrossClusterSubjectReview{
		multiClusterClient: mClient,
		scheme:             scheme,
		restMapper:         restMapper,
		ClusterParameter:   "cluster",
		ClusterNamespace:   system.Namespace(),
	}
}

// CrossClusterSubjectReview describe a struct to get the client of special cluster and simulate the requesting user
type CrossClusterSubjectReview struct {
	multiClusterClient Interface
	scheme             *runtime.Scheme
	restMapper         meta.RESTMapper

	ClusterParameter string
	ClusterNamespace string
}

// SetClusterParameter sets the cluster parameter name
func (c *CrossClusterSubjectReview) SetClusterParameter(parameter string) {
	c.ClusterParameter = parameter
}

// SetClusterNamespace set the namespace which the cluster resource is stored in
func (c *CrossClusterSubjectReview) SetClusterNamespace(ns string) {
	c.ClusterNamespace = ns
}

func (c *CrossClusterSubjectReview) getClusterParameterName(req *restful.Request) string {
	if name := req.PathParameter(c.ClusterParameter); name != "" {
		return name
	}
	return req.QueryParameter(c.ClusterParameter)
}

// GetClient get k8s client of the specified cluster and simulate the requesting user
func (c *CrossClusterSubjectReview) GetClient(ctx context.Context, req *restful.Request) (client.Client, error) {
	clusterName := c.getClusterParameterName(req)
	if clusterName == "" {
		return nil, nil
	}
	clusterRef := &corev1.ObjectReference{}
	clusterRef.SetGroupVersionKind(ClusterRegistryGVK)
	clusterRef.Name = clusterName
	clusterRef.Namespace = c.ClusterNamespace

	config, err := c.multiClusterClient.GetConfig(ctx, clusterRef)
	if err != nil {
		return nil, err
	}
	reqCtx := req.Request.Context()
	u, ok := apiserverrequest.UserFrom(reqCtx)
	if !ok {
		return nil, nil
	}
	copyConfig := rest.CopyConfig(config)
	copyConfig.Impersonate.UID = u.GetUID()
	copyConfig.Impersonate.Groups = u.GetGroups()
	copyConfig.Impersonate.UserName = u.GetName()
	copyConfig.Impersonate.Extra = u.GetExtra()
	directClient, err := client.New(copyConfig, client.Options{Scheme: c.scheme, Mapper: c.restMapper})
	if err != nil {
		return nil, err
	}
	return directClient, nil
}
