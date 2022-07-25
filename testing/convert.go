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

package testing

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConvertCorev1Resouces convert runtime.Object to client.Object
// be similar to convertFromUnstructuredIfNecessary
func ConvertCorev1Resouces(runtimeObj runtime.Object) (obj client.Object, err error) {
	switch v := runtimeObj.(type) {
	case *corev1.PersistentVolume:
		obj = v
	case *corev1.PersistentVolumeClaim:
		obj = v
	case *corev1.Pod:
		obj = v
	case *corev1.ReplicationController:
		obj = v
	case *corev1.Service:
		obj = v
	case *corev1.ServiceAccount:
		obj = v
	case *corev1.Endpoints:
		obj = v
	case *corev1.Node:
		obj = v
	case *corev1.Namespace:
		obj = v
	case *corev1.Binding:
		obj = v
	case *corev1.LimitRange:
		obj = v
	case *corev1.ResourceQuota:
		obj = v
	case *corev1.Secret:
		obj = v
	case *corev1.ConfigMap:
		obj = v
	default:
		err = fmt.Errorf("Unsupported gvk: %s", runtimeObj.GetObjectKind().GroupVersionKind())
	}
	return
}
