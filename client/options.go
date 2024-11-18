/*
Copyright 2023 The AlaudaDevops Authors.

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

package client

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetOptions is a wrapper for client.GetOptions
type GetOptions struct {
	client.GetOptions
}

// NewGetOptions returns a new GetOptions
func NewGetOptions() *GetOptions {
	return &GetOptions{}
}

// init returns a new GetOptions if opt is nil or opt.Raw is nil
func (opt *GetOptions) init() *GetOptions {
	if opt != nil && opt.GetOptions.Raw != nil {
		return opt
	}
	return &GetOptions{
		GetOptions: client.GetOptions{
			Raw: &metav1.GetOptions{},
		},
	}
}

// WithCached set the ResourceVersion to 0
func (opt *GetOptions) WithCached() *GetOptions {
	opt = opt.init()
	opt.Raw.ResourceVersion = "0"
	return opt
}

// Build returns the client.GetOptions
func (opt *GetOptions) Build() *client.GetOptions {
	return &opt.GetOptions
}

// GetCachedOptions returns GetOptions with ResourceVersion set to 0
func GetCachedOptions() *client.GetOptions {
	return NewGetOptions().WithCached().Build()
}

// ListOptions is a wrapper for client.ListOptions
type ListOptions struct {
	client.ListOptions
}

// NewListOptions returns a new ListOptions
func NewListOptions() *ListOptions {
	return &ListOptions{}
}

// init returns a new ListOptions if opt is nil or opt.Raw is nil
func (opt *ListOptions) init() *ListOptions {
	if opt != nil && opt.ListOptions.Raw != nil {
		return opt
	}
	return &ListOptions{
		ListOptions: client.ListOptions{
			Raw: &metav1.ListOptions{},
		},
	}
}

// WithCached set the ResourceVersion to 0
func (opt *ListOptions) WithCached() *ListOptions {
	opt = opt.init()
	opt.Raw.ResourceVersion = "0"
	return opt
}

// WithLimit set the limit
func (opt *ListOptions) WithLimit(limit int64) *ListOptions {
	opt = opt.init()
	opt.Limit = limit
	return opt
}

// WithNamespace set the namespace
func (opt *ListOptions) WithNamespace(namespace string) *ListOptions {
	opt = opt.init()
	opt.Namespace = namespace
	return opt
}

// WithUnsafeDisableDeepCopy set the UnsafeDisableDeepCopy to true
func (opt *ListOptions) WithUnsafeDisableDeepCopy() *ListOptions {
	opt = opt.init()
	True := true
	opt.UnsafeDisableDeepCopy = &True
	return opt
}

// Build returns the client.ListOptions
func (opt *ListOptions) Build() *client.ListOptions {
	return &opt.ListOptions
}

// ListCachedOptions returns ListOptions with ResourceVersion set to 0
func ListCachedOptions() *client.ListOptions {
	return NewListOptions().WithCached().Build()
}
