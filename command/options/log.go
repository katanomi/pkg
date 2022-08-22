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

	"github.com/katanomi/pkg/command/io"
	"github.com/katanomi/pkg/command/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log Log related options
type Log struct {
	verbose bool
	Logger  *zap.SugaredLogger
}

// Enabled decides whether a given logging level is enabled
func (opts *Log) Enabled(l zapcore.Level) bool {
	if opts.verbose {
		return true
	}

	return l >= zapcore.InfoLevel
}

// Setup set up the Log
func (opts *Log) Setup(ctx context.Context, cmd *cobra.Command, args []string) {
	if opts.Logger == nil {
		if l := logger.GetLogger(ctx); l != nil {
			opts.Logger = l
		} else {
			iostreams := io.MustGetIOStreams(ctx)
			opts.Logger = logger.NewLogger(zapcore.AddSync(iostreams.ErrOut), opts)
		}
	}
}

// AddFlags add flags to options
func (opts *Log) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&opts.verbose, `verbose`, `v`, false, `sets the Log level to be displayed.`)
}
