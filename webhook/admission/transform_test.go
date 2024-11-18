/*
Copyright 2023 The AlaudaDevops Authors.

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

package admission

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlaudaDevops/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestWithUpdateTime(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	obj := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
			UID:       types.UID("abc"),
		},
	}

	req := admission.Request{
		AdmissionRequest: v1.AdmissionRequest{
			Operation: v1.Create,
		},
	}
	WithUpdateTime()(ctx, obj, req)
	g.Expect(obj.Annotations).NotTo(HaveKey(v1alpha1.UpdatedTimeAnnotationKey))

	req.Operation = v1.Update
	WithUpdateTime()(ctx, obj, req)
	beTimeString := matchers.NewWithTransformMatcher(checkTimeString, BeTrue())
	g.Expect(obj.Annotations[v1alpha1.UpdatedTimeAnnotationKey]).To(beTimeString)
}

func checkTimeString(actual interface{}) (interface{}, error) {
	str, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("expected string value, got %v", actual)
	}
	_, err := time.Parse(time.RFC3339, str)
	return err == nil, err
}
