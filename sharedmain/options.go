/*
Copyright 2021 The AlaudaDevops Authors.

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

package sharedmain

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrlcluster "sigs.k8s.io/controller-runtime/pkg/cluster"
)

// WithScheme adds a scheme to cluster.Options
func WithScheme(scheme *runtime.Scheme) ctrlcluster.Option {
	return func(opts *ctrlcluster.Options) {
		opts.Scheme = scheme
	}
}
