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

package resource

import (
	corev1 "k8s.io/api/core/v1"
)

// SumResources sums many ResourceRequirements resources.
// Note: No operations are performed on Claims. only the sum of Limits and Requests is completed.
func SumResources(resources ...corev1.ResourceRequirements) corev1.ResourceRequirements {
	if len(resources) == 0 {
		return corev1.ResourceRequirements{}
	}

	sumLimits := make([]corev1.ResourceList, 0, len(resources))
	sumRequests := make([]corev1.ResourceList, 0, len(resources))
	for _, item := range resources {
		sumLimits = append(sumLimits, item.Limits)
		sumRequests = append(sumRequests, item.Requests)
	}

	return corev1.ResourceRequirements{
		Limits:   SumResourceList(sumLimits...),
		Requests: SumResourceList(sumRequests...),
	}
}

// SumResourceList sums many ResourceList resources.
func SumResourceList(list ...corev1.ResourceList) corev1.ResourceList {
	result := corev1.ResourceList{}
	for _, item := range list {
		for key, itemQuantity := range item {
			tmpQuantity := itemQuantity
			if v, ok := result[key]; ok {
				v.Add(itemQuantity)
				result[key] = v
			} else {
				result[key] = tmpQuantity
			}
		}
	}

	return result
}
