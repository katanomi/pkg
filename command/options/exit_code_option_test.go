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

package options

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

func TestExitCodeOption(t *testing.T) {
	t.Run("exit code is succeed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			ExitCodeOption
		}{}
		args := []string{
			"--exit-code-path", "./testdata/exit-code-option-succeed.txt",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		success, code, err := obj.Succeed()
		g.Expect(err).To(Succeed())
		g.Expect(code).To(Equal("0"))
		g.Expect(success).To(Equal(true))
	})

	t.Run("exit code file not found", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			ExitCodeOption
		}{}
		args := []string{
			"--exit-code-path", "./testdata/not-found-path.txt",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		success, code, err := obj.Succeed()
		g.Expect(err).ToNot(Succeed())
		g.Expect(code).To(Equal(""))
		g.Expect(success).To(Equal(false))
	})

	t.Run("exit code file not set", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			ExitCodeOption
		}{}
		args := []string{}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		success, code, err := obj.Succeed()
		g.Expect(err).ToNot(Succeed())
		g.Expect(code).To(Equal(""))
		g.Expect(success).To(Equal(false))
	})

	t.Run("exit code non-zero exit code", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			ExitCodeOption
		}{}
		args := []string{
			"--exit-code-path", "./testdata/exit-code-option-failed.txt",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		success, code, err := obj.Succeed()
		g.Expect(err).To(Succeed())
		g.Expect(code).To(Equal("-1"))
		g.Expect(success).To(Equal(false))
	})

	t.Run("exit code file context is empty", func(t *testing.T) {
		g := NewGomegaWithT(t)

		obj := struct {
			ExitCodeOption
		}{}
		args := []string{
			"--exit-code-path", "./testdata/exit-code-option-empty.txt",
		}

		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterFlags(&obj, flagSet)
		err := flagSet.Parse(args)
		g.Expect(err).Should(Succeed(), "parse flag succeed.")
		success, code, err := obj.Succeed()
		g.Expect(err).To(Succeed())
		g.Expect(code).To(Equal(""))
		g.Expect(success).To(Equal(true))
	})
}
