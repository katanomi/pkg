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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CopyLabels from the left side object to the righ side
// will override any existing labels
func CopyLabels(object, dest metav1.Object) {
	labels := dest.GetLabels()
	originalLabels := object.GetLabels()
	labels = CopyMapStringString(originalLabels, labels)
	dest.SetLabels(labels)
}

// CopyAnnotations from the left side object to the righ side
// will override any existing values
func CopyAnnotations(object, dest metav1.Object) {
	anno := dest.GetAnnotations()
	originalAnno := object.GetAnnotations()
	anno = CopyMapStringString(originalAnno, anno)
	dest.SetAnnotations(anno)
}

// CopyMapStringString copies content from a map to annother
func CopyMapStringString(object, dest map[string]string) map[string]string {
	if object != nil {
		if dest == nil {
			dest = map[string]string{}
		}
		for k, v := range object {
			dest[k] = v
		}
	}
	return dest
}
