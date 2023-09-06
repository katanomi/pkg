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

package secret

import (
	"context"
	"errors"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// MockMutator is a mock implementation of the Mutator interface for testing.
type MockMutator struct {
	Modified bool
}

// Mutate modifies the Secret object by appending "-modified" to its Name field.
func (m *MockMutator) Mutate(_ context.Context, secret *corev1.Secret) error {
	secret.Name += "-modified"
	m.Modified = true
	return nil
}

// TestMutatorList tests the Mutate function for a list of Mutators.
func TestMutatorList(t *testing.T) {
	g := NewGomegaWithT(t)

	mockMutator := &MockMutator{}
	list := MutatorList{DefaultMutator, mockMutator}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	err := list.Mutate(context.Background(), secret)

	g.Expect(err).To(BeNil())
	g.Expect(secret.Name).To(Equal("test-modified"))
	g.Expect(mockMutator.Modified).To(BeTrue())
}

// TestMutatorList tests the Mutate function for a list of Mutators.
func TestMutatorListError(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := MutateFunc(func(ctx context.Context, secret *corev1.Secret) error {
		return errors.New("test")
	})

	m2 := MutateFunc(func(ctx context.Context, secret *corev1.Secret) error {
		return nil
	})
	list := MutatorList{m1, m2}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	err := list.Mutate(context.Background(), secret)

	g.Expect(err).NotTo(BeNil())
}

// TestNoOpMutator tests the no-operation Mutator.
func TestNoOpMutator(t *testing.T) {
	g := NewGomegaWithT(t)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	err := DefaultMutator.Mutate(context.Background(), secret)

	g.Expect(err).To(BeNil())
	g.Expect(secret.Name).To(Equal("test"))
}

// TestNoOpMutator tests the no-operation Mutator.
func TestMutatorFunc(t *testing.T) {
	g := NewGomegaWithT(t)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	mutator := MutateFunc(func(ctx context.Context, secret *corev1.Secret) error {
		secret.Name += "-modified"
		return nil
	})
	err := mutator.Mutate(context.Background(), secret)

	g.Expect(err).To(BeNil())
	g.Expect(secret.Name).To(Equal("test-modified"))
}

// TestMutatorFromCtx tests the retrieval of a Mutator from a context.
func TestMutatorFromCtx(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	g.Expect(MutatorFromCtx(ctx)).To(Equal(DefaultMutator))

	mockMutator := &MockMutator{}
	ctx = WithMutator(ctx, mockMutator)
	g.Expect(MutatorFromCtx(ctx)).To(Equal(mockMutator))
}
