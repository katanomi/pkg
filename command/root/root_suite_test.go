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

package root_test

import (
	// "bytes"
	"testing"

	"context"

	"github.com/katanomi/pkg/command/io"
	"github.com/katanomi/pkg/command/root"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	clioptions "k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestRoot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Root Suite")
}

var _ = Describe("NewRootCommand", func() {
	var (
		ctx         context.Context
		streams     clioptions.IOStreams
		cmd         *cobra.Command
		subcommands []root.SubcommandFunc
		err         error
		// out  *bytes.Buffer
	)

	BeforeEach(func() {
		streams, _, _, _ = clioptions.NewTestIOStreams()
		// uses ginkgo writer for logs
		streams.ErrOut = GinkgoWriter
		ctx = context.Background()
		ctx = io.WithIOStreams(ctx, &streams)
		subcommands = nil
	})

	JustBeforeEach(func() {
		cmd = root.NewRootCommand(ctx, "test-cli", subcommands...)
		err = cmd.Execute()

	})
	It("cmd is not nil", func() {
		Expect(cmd).ToNot(BeNil())
	})

	When("with subcommands", func() {
		BeforeEach(func() {
			subcommands = append(subcommands, func(_ context.Context, _ string) *cobra.Command {
				return &cobra.Command{Use: "subcommamd", Run: func(_ *cobra.Command, _ []string) {}}
			})
		})
		It("should have subcommand", func() {
			// by default adds completion and help subcommands
			Expect(cmd.Commands()).To(HaveLen(3), "should have subcommands")
			Expect(err).To(BeNil(), "should not error")
		})

	})
	When("without subcommands", func() {
		It("should NOT have subcommands", func() {
			Expect(cmd.Commands()).To(HaveLen(0), "should NOT have subcommands")
			Expect(err).To(BeNil(), "should not error")
		})
	})
})
