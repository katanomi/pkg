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

import corev1 "k8s.io/api/core/v1"

// SumResources sums two ResourceRequirements resources.
// Note: No operations are performed on Claims. only the sum of Limits and Requests is completed.
func SumResources(param1, param2 corev1.ResourceRequirements) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits:   SumResourceList(param1.Limits, param2.Limits),
		Requests: SumResourceList(param1.Requests, param2.Requests),
	}
}

// SumResourceList sums two ResourceList resources.
func SumResourceList(param1, param2 corev1.ResourceList) corev1.ResourceList {
	result := corev1.ResourceList{}
	for key, itemQuantity := range param1 {
		tmpQuantity := itemQuantity
		if quantity, ok := param2[key]; ok {
			tmpQuantity.Add(quantity)
			result[key] = tmpQuantity
			delete(param2, key)
		} else {
			result[key] = tmpQuantity
		}
	}

	for key, itemQuantity := range param2 {
		tmpQuantity := itemQuantity
		result[key] = tmpQuantity
	}
	return result
}
