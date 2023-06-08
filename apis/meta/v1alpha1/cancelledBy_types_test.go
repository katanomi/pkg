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

	. "github.com/onsi/gomega"
)

func TestCancellerAnnotation(t *testing.T) {
	RegisterTestingT(t)

	cancelledBy := CancelledBy{}

	t.Run("FromAnnotationCanceller", func(t *testing.T) {
		annotations := map[string]string{
			CancelledByAnnotationKey: `{"user":{"kind":"User","name":"admin"},"ref":{"kind":"BuildRun","namespace":"default","name":"test"}}`,
		}
		cancelledBy, err := cancelledBy.FromAnnotationCancelledBy(annotations)
		Expect(err).To(BeNil())
		Expect(cancelledBy.Ref.Name).To(Equal("test"))
		Expect(cancelledBy.User.Kind).To(Equal("User"))
		Expect(cancelledBy.User.Name).To(Equal("admin"))
	})

	t.Run("FromAnnotationCanceller with invalid data", func(t *testing.T) {
		annotations := map[string]string{
			CancelledByAnnotationKey: `{"user":{"kind"namespace":"default","name":"test"}}`,
		}
		_, err := cancelledBy.FromAnnotationCancelledBy(annotations)
		Expect(err).ToNot(BeNil())
	})
}

func TestSetIntoAnnotationCancelledBy(t *testing.T) {
	RegisterTestingT(t)

	cancelledBy := CancelledBy{}

	t.Run("FromAnnotationCanceller", func(t *testing.T) {
		annotations, err := cancelledBy.SetIntoAnnotationCancelledBy(nil)
		Expect(err).To(BeNil())
		Expect(annotations).ToNot(BeNil())
	})
}
