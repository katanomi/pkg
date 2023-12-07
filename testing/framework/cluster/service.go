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

package cluster

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PluginService Get plugin service by plugin name
func PluginService(ctx context.Context, clt client.Client, name string) (*corev1.Service, error) {
	services := &corev1.ServiceList{}
	if err := clt.List(ctx, services, client.MatchingLabels{"plugin": name}); err != nil {
		return nil, err
	}

	if len(services.Items) > 0 {
		return &services.Items[0], nil
	}
	return nil, nil
}
