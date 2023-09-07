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

	corev1 "k8s.io/api/core/v1"
)

// Mutator interface defines a method for modifying Secret objects.
type Mutator interface {
	Mutate(context.Context, *corev1.Secret) error
}

// MutateFunc is a function type that implements the Mutator interface.
type MutateFunc func(context.Context, *corev1.Secret) error

func (m MutateFunc) Mutate(ctx context.Context, secret *corev1.Secret) error {
	return m(ctx, secret)
}

// MutatorList is a slice of Mutators.
type MutatorList []Mutator

// Mutate iterates over the list of mutators and applies them to the provided Secret object.
func (m MutatorList) Mutate(ctx context.Context, secret *corev1.Secret) error {
	for _, mutator := range m {
		if err := mutator.Mutate(ctx, secret); err != nil {
			return err
		}
	}
	return nil
}

// DefaultMutator is a no-operation implementation of the Mutator interface.
// It is used as a default option when there is no Mutator associated with the context.
var DefaultMutator = &NoOpMutator{}

// NoOpMutator struct implements the Mutator interface but performs no actions.
type NoOpMutator struct {
}

// Mutate is the NoOpMutator's implementation that simply returns the passed-in Secret object without any modifications.
func (a *NoOpMutator) Mutate(_ context.Context, secret *corev1.Secret) error {
	return nil
}

// mutatorKey is a key used for storing a Mutator in a context.
var mutatorKey = struct{}{}

// MutatorFromCtx retrieves a Mutator from the given context.
// If no Mutator is found, it returns the noOpMutator.
func MutatorFromCtx(ctx context.Context) Mutator {
	if mutation, ok := ctx.Value(mutatorKey).(Mutator); ok {
		return mutation
	}
	return DefaultMutator
}

// WithMutator attaches a Mutator to a context and returns the new context.
func WithMutator(ctx context.Context, mutator Mutator) context.Context {
	return context.WithValue(ctx, mutatorKey, mutator)
}
