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

package exec

import (
	"bytes"
	"context"
	"io"
)

// NoOpCmder only writes down all executed commands as output
// used for testing only
type NoOpCmder struct{}

// Command initializes a NoOpCmd with given arguments
func (n *NoOpCmder) Command(cmd string, args ...string) Cmd {
	return &NoOpCmd{Command: cmd, Args: args}
}

// CommandContext initializes a NoOpCmd with given arguments
func (n *NoOpCmder) CommandContext(ctx context.Context, cmd string, args ...string) Cmd {
	return &NoOpCmd{Context: ctx, Command: cmd, Args: args}
}

// NoOpCmd only writes down all executed commands as output
// used for testing only
type NoOpCmd struct {
	Context context.Context
	Command string
	Args    []string
	Envs    []string

	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// Run writes down its command and args to the Stdout writer
func (no *NoOpCmd) Run() error {
	buff := &bytes.Buffer{}
	buff.WriteString(no.Command)
	buff.WriteString(" ")
	for i, arg := range no.Args {
		buff.WriteString(arg)
		if i+1 < len(no.Args) {
			buff.WriteString(" ")
		}
	}
	_, err := no.Stdout.Write(buff.Bytes())
	return err
}

// SetEnv sets env
func (no *NoOpCmd) SetEnv(envs ...string) Cmd {
	no.Envs = envs
	return no
}

// SetStdin sets Stdin reader, will not be used
func (no *NoOpCmd) SetStdin(stdin io.Reader) Cmd {
	no.Stdin = stdin
	return no
}

// SetStdout sets Stdout
func (no *NoOpCmd) SetStdout(stdout io.Writer) Cmd {
	no.Stdout = stdout
	return no
}

// SetStdout sets Stderr
func (no *NoOpCmd) SetStderr(stderr io.Writer) Cmd {
	no.Stderr = stderr
	return no
}
