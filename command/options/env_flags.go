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
	"context"

	"github.com/joho/godotenv"
	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/spf13/cobra"
)

type envValuesValidationKey struct{}

// EnvFlagsOption describe env flags option
type EnvFlagsOption struct {
	EnvFlags map[string]string
}

// ValuesValidationOptionsFrom returns the value of the envValuesValidationKey key on the ctx
func ValuesValidationOptionsFrom(ctx context.Context) ([]pkgargs.ValuesValidateOption, bool) {
	opts, ok := ctx.Value(envValuesValidationKey{}).([]pkgargs.ValuesValidateOption)
	return opts, ok
}

// WithValuesValidationOpts returns a copy of parent in which the []ValuesValidateOption is set
func WithValuesValidationOpts(parent context.Context, opts []pkgargs.ValuesValidateOption) context.Context {
	return context.WithValue(parent, envValuesValidationKey{}, opts)
}

// Setup defines how to start with env-flags option
func (p *EnvFlagsOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	opts, ok := ValuesValidationOptionsFrom(ctx)
	if !ok {
		opts = []pkgargs.ValuesValidateOption{pkgargs.ValuesValidationOptDuplicatedKeys}
	}
	p.EnvFlags, err = pkgargs.GetKeyValues(ctx, args, "env-flags", opts...)
	if err != nil {
		return err
	}
	return nil
}

// WriteEnvFile writes env variables to target file
func (p *EnvFlagsOption) WriteEnvFile(filename string) error {
	return godotenv.Write(p.EnvFlags, filename)
}
