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

package root

import (
	"context"
	"fmt"

	"github.com/katanomi/pkg/command/io"
	"github.com/katanomi/pkg/command/logger"
	"github.com/katanomi/pkg/command/options"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

// SubcommandFunc inits a subcommand to be inserted inside root
type SubcommandFunc func(ctx context.Context, name string) *cobra.Command

// NewRootCommand initiates all commands. This is the main entrypoint of the cli
func NewRootCommand(ctx context.Context, name string, subcommands ...SubcommandFunc) *cobra.Command {
	// sets log as persistent options and provides logger using
	// context variables
	logOpts := &options.Log{}
	streams := io.MustGetIOStreams(ctx)
	ctx = logger.WithLogger(ctx, logger.NewLogger(zapcore.AddSync(streams.ErrOut), logOpts))
	rootCmd := &cobra.Command{
		Use: fmt.Sprintf("%s [command] [options]", name),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logOpts.Setup(ctx, cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	// will persist flag across all subcommands
	logOpts.AddFlags(rootCmd.PersistentFlags())

	for _, sub := range subcommands {
		rootCmd.AddCommand(sub(ctx, name))
	}

	return rootCmd
}
