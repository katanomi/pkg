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
	"strings"

	pkgargs "github.com/katanomi/pkg/command/args"

	"github.com/spf13/cobra"
)

// ScanFlagsOption describe scan flags option
type ScanFlagsOption struct {
	ScanFlags map[string]string
}

func (p *ScanFlagsOption) Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	p.ScanFlags = make(map[string]string)
	scanFlags, _ := pkgargs.GetKeyValues(ctx, args, "scan-flags")

	for k, v := range scanFlags {
		p.ScanFlags[strings.Trim(k, ".")] = v
	}
	return nil
}
