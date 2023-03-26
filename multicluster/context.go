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

// Package multicluster contains functions to add and retrieve multi cluster from context
package multicluster

import (
	"context"
)

type multiClusterKey struct{}

// WithMultiCluster adds a multi cluster client to the context
func WithMultiCluster(ctx context.Context, clt Interface) context.Context {
	return context.WithValue(ctx, multiClusterKey{}, clt)
}

// MultiCluster returns a multicluster client in context
func MultiCluster(ctx context.Context) Interface {
	val := ctx.Value(multiClusterKey{})
	if val == nil {
		return nil
	}
	clt, _ := val.(Interface)
	return clt
}

type clusterNamesKey struct{}

// WithClusterNames adds cluster names to the context
func WithClusterNames(ctx context.Context, names []string) context.Context {
	return context.WithValue(ctx, clusterNamesKey{}, names)
}

// ClusterNames return a cluster name list in context
func ClusterNames(ctx context.Context) []string {
	val := ctx.Value(clusterNamesKey{})
	if val == nil {
		return nil
	}
	names, _ := val.([]string)
	return names
}

type ignoreForbiddenKey struct{}

// WithIgnoreForbidden adds ignore forbidden flag to the context
func WithIgnoreForbidden(ctx context.Context, ignoreForbidden bool) context.Context {
	return context.WithValue(ctx, ignoreForbiddenKey{}, ignoreForbidden)
}

// IgnoreForbidden return a ignore forbidden flag in context
func IgnoreForbidden(ctx context.Context) bool {
	val := ctx.Value(ignoreForbiddenKey{})
	if val == nil {
		return false
	}
	ignoreForbidden, _ := val.(bool)
	return ignoreForbidden
}
