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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestGetNamespacedNameFromObject(t *testing.T) {
	table := map[string]struct {
		Object metav1.Object
		Result types.NamespacedName
	}{
		"Simple secret object": {
			Object: &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
				Name: "secret", Namespace: "default",
			}},
			Result: types.NamespacedName{Name: "secret", Namespace: "default"},
		},
		"Nil object": {
			Object: nil,
			Result: types.NamespacedName{},
		},
	}

	for name, item := range table {
		test := item
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := GetNamespacedNameFromObject(test.Object)
			diff := cmp.Diff(test.Result, result)

			g.Expect(diff).To(BeEmpty())
		})
	}
}
