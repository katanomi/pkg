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
	"testing"

	"github.com/AlaudaDevops/pkg/fieldindexer"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAppWithFieldIndexer(t *testing.T) {
	t.Run("append one field indexer ", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := App("test").WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(1))
	})
	t.Run("append more than one field indexer", func(t *testing.T) {
		g := NewGomegaWithT(t)

		a := App("test").WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		}).WithFieldIndexer(fieldindexer.FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.name",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["name"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(2))
	})
}
