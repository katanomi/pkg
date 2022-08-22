/*
Copyright 2022 The Katanomi Authors.

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

package io

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	clioptions "k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestWithIOStreams(t *testing.T) {
	var (
		ctx       context.Context
		result    context.Context
		iostreams *clioptions.IOStreams
	)
	t.Run("nil context", func(t *testing.T) {
		g := NewGomegaWithT(t)
		result = WithIOStreams(ctx, iostreams)
		g.Expect(result).ToNot(BeNil())

		iostreams = GetIOStreams(ctx)
		// expect an empty struct if not stored previously
		g.Expect(iostreams).To(BeNil())

		iostreams = MustGetIOStreams(ctx)
		g.Expect(iostreams).ToNot(BeNil())
		g.Expect(*iostreams).To(Equal(clioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		}))
	})
	t.Run("check if is the same iostreams", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctx = context.Background()
		discardStream := clioptions.NewTestIOStreamsDiscard()
		iostreams = &discardStream
		ctx = WithIOStreams(ctx, iostreams)
		g.Expect(result).ToNot(BeNil())

		storedStreams := GetIOStreams(ctx)
		g.Expect(storedStreams).To(Equal(iostreams))

		storedStreams = MustGetIOStreams(ctx)
		g.Expect(storedStreams).To(Equal(iostreams))
	})
}
