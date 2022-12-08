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
	"fmt"

	"github.com/joho/godotenv"
	pkgargs "github.com/katanomi/pkg/command/args"
	"github.com/spf13/cobra"
)

// EnvFlagsOption describe env flags option
type EnvFlagsOption struct {
	EnvFlags map[string]string
}

// Setup defines how to start with env-flags option
func (p *EnvFlagsOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	p.EnvFlags, _ = pkgargs.GetKeyValues(ctx, args, "env-flags")
	pairs, _ := pkgargs.GetArrayValues(ctx, args, "env-flags")
	if len(pairs) != len(p.EnvFlags) {
		return fmt.Errorf("invalid env-flags")
	}
	return nil
}

// WriteEnvFile writes env variables to target file
func (p *EnvFlagsOption) WriteEnvFile(filename string) error {
	return godotenv.Write(p.EnvFlags, filename)
}
