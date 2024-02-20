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

package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers"
	v1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestSetCancelledBy(t *testing.T) {
	RegisterTestingT(t)

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
			UserInfo: authenticationv1.UserInfo{
				Username: "admin@xxxx.io",
			},
		},
	}
	t.Run("setCancelledBy", func(t *testing.T) {
		setCancelledBy(ctx, obj, req)
		Expect(len(obj.Annotations)).To(Equal(1))
		Expect(obj.Annotations[v1alpha1.CancelledByAnnotationKey]).To(Equal(`{"user":{"kind":"User","name":"admin@xxxx.io"}}`))
	})

}

func TestWithCancelledBy(t *testing.T) {
	RegisterTestingT(t)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	result, err := json.Marshal(pod)
	Expect(err).To(BeNil())

	schema := runtime.NewScheme()
	corev1.AddToScheme(schema)
	cancelledBy := WithCancelledBy(schema, func(runtime.Object, runtime.Object) bool { return true })
	request := admission.Request{
		AdmissionRequest: v1.AdmissionRequest{
			UserInfo: authenticationv1.UserInfo{
				Username: "admin",
			},
			Operation: v1.Update,
			OldObject: runtime.RawExtension{
				Raw: result,
			},
		},
	}
	ctx := context.Background()

	cancelledBy(ctx, pod, request)
	Expect(pod.Annotations[v1alpha1.CancelledByAnnotationKey]).To(Equal(`{"user":{"kind":"User","name":"admin"}}`))

}

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
