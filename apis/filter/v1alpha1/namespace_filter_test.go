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

package filter

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNamespaceFilterRule_Filter(t *testing.T) {
	namespaces := []corev1.Namespace{
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"test-1-k": "test-1",
				},
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-2",
				Annotations: map[string]string{
					"test-2-k": "test-2",
				},
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-3",
				Labels: map[string]string{
					"test-3/k": "test-3",
				},
				Annotations: map[string]string{
					"test-3/k": "test-3",
				},
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-4",
				Labels: map[string]string{
					"test-4/l": "test-4",
				},
				Annotations: map[string]string{
					"test-4/a": "test-4",
				},
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-5",
				Labels: map[string]string{
					"test-5.label": "test-5",
				},
				Annotations: map[string]string{
					"test-5.annotation": "test-5",
				},
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-6",
				Labels: map[string]string{
					"test-6.label": "test-6",
				},
				Annotations: map[string]string{
					"test-6.annotation": "test-6",
				},
			},
		},
	}

	cases := []struct {
		Exact    map[string]string
		Expected int
	}{
		{
			Exact: map[string]string{
				"$(metadata.labels.test-1-k)": "$(metadata.name)",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(metadata.annotations.test-2-k)": "$(metadata.name)",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(metadata.labels.test-3/k)":      "$(metadata.name)",
				"$(metadata.annotations.test-3/k)": "$(metadata.name)",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(metadata.labels.test-4/l)": "$(metadata.annotations.test-4/a)",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				`$(metadata.labels.test-5\.label)`:           "$(metadata.name)",
				`$(metadata.annotations.test-5\.annotation)`: "$(metadata.name)",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(metadata.labels.test-1-k)": "test-1",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(metadata.name)": "test-1",
			},
			Expected: 1,
		},
		{
			Exact: map[string]string{
				"$(some-key-1)": "$(some-key-2)",
			},
			Expected: 0,
		},
		{
			Exact: map[string]string{
				"$(metadata.name)qwer": "test-1",
			},
			Expected: 0,
		},
		{
			Exact: map[string]string{
				`$(metadata.labels.test-6\\.label)`:           "$(metadata.name)",
				`$(metadata.annotations.test-6\\.annotation)`: "$(metadata.name)",
			},
			Expected: 1,
		},
	}

	for i := range cases {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			g := NewGomegaWithT(t)

			rule := NamespaceFilterRule{Exact: cases[i].Exact}

			g.Expect(rule.Filter(namespaces)).To(HaveLen(cases[i].Expected))
		})
	}
}
