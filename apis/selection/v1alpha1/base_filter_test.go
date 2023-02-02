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

package v1alpha1

import (
	"testing"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("TestBaseFilter_MatchObject", func() {
	var (
		filter   *BaseFilter
		cm       *corev1.ConfigMap
		matchRet bool
	)
	BeforeEach(func() {
		matchRet = false
		filter = &BaseFilter{}
		cm = &corev1.ConfigMap{}
		ktesting.MustLoadYaml("./testdata/baseFilter.configmap.yaml", cm)
	})
	JustBeforeEach(func() {
		matchRet = filter.MatchObject(cm)
	})
	Context("match succeed", func() {
		When("matched by selector", func() {
			BeforeEach(func() {
				filter.Selector = &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				}
			})
			It("should return true", func() {
				Expect(matchRet).To(BeTrue())
			})
		})
		When("matched by refs", func() {
			BeforeEach(func() {
				ref := metav1alpha1.GetObjectReferenceFromObject(cm,
					metav1alpha1.ObjectRefWithNamespace(),
					metav1alpha1.ObjectRefWithTypeMeta())
				filter.Refs = []corev1.ObjectReference{ref}
			})
			It("should return true", func() {
				Expect(matchRet).To(BeTrue())
			})
		})
		When("matched by filter", func() {
			BeforeEach(func() {
				filter.Selector = &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				}
				filter.Filter = &BaseFilterRule{
					Exact: map[string]string{
						"$(metadata.name)": "not-exist-name",
					},
				}
			})
			It("should return false", func() {
				Expect(matchRet).To(BeFalse())
			})
		})
	})
	Context("match failed", func() {
		It("should return false", func() {
			Expect(matchRet).To(BeFalse())
		})
	})
})

func TestBaseFilter_MatchObject(t *testing.T) {
	type fields struct {
		Selector *metav1.LabelSelector
		Filter   *BaseFilterRule
		Refs     []corev1.ObjectReference
	}
	type args struct {
		obj client.Object
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BaseFilter{
				Selector: tt.fields.Selector,
				Filter:   tt.fields.Filter,
				Refs:     tt.fields.Refs,
			}
			if got := p.MatchObject(tt.args.obj); got != tt.want {
				t.Errorf("MatchObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
