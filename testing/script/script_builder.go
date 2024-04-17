/*
Copyright 2024 The Katanomi Authors.

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

package script

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	testingcontext "github.com/katanomi/pkg/testing/context"
	"knative.dev/pkg/logging"
)

// ScriptBuilder script runner for easier initializing of data using scripts
// while providing option for asserting and processing the result
type ScriptBuilder struct {
	ctx    context.Context
	script string
	args   []string

	repoPath string

	rootFolder string

	output       io.Writer
	outputPrefix string

	errorOutput io.Writer
}

// NewScript starts creating a script runner for the given script and args
// which can be customized by chaining other methods
// and be finally executed by calling the Result or Error methods
func NewScript(ctx context.Context, scriptName string, args ...string) *ScriptBuilder {
	builder := &ScriptBuilder{
		ctx:        ctx,
		rootFolder: e2eTempDir,
		script:     scriptName,
		args:       args,
	}
	return builder.WithLocalRepoPath()
}

// WithRepoPath changes the repository path.
// By default this value is already initiated using WithLocalRepoPath
// which will fetch a LocalRepoPath from the context
func (r *ScriptBuilder) WithRepoPath(repoPath string) *ScriptBuilder {
	r.repoPath = repoPath
	return r
}

// WithLocalRepoPath set the repo path from the LocalRepoPathFromCtx
// function, which normally is already present inside context when
// adding repo conditions like
func (r *ScriptBuilder) WithLocalRepoPath() *ScriptBuilder {
	localRepoPath := testingcontext.LocalRepoPathFromCtx(r.ctx)
	if localRepoPath != nil {
		r.repoPath = *localRepoPath
	}
	return r
}

// WithRootFolder changes the root folder of the scripts.
// By default this value is already initiated using a
// variable that can be customized using the E2E_TEMP_DIR environment variable
func (r *ScriptBuilder) WithRootFolder(rootFolder string) *ScriptBuilder {
	r.rootFolder = rootFolder
	return r
}

// WithScript changes the script initialized by the builder
func (r *ScriptBuilder) WithScript(scriptName string, args ...string) *ScriptBuilder {
	r.script = scriptName
	r.args = args
	return r
}

// WithOutput adds an output writer to write text output from
// script into. The second argument accepts an output prefix from the script
// which defaults to "##output##
func (r *ScriptBuilder) WithOutput(writer io.Writer, prefix ...string) *ScriptBuilder {
	r.output = writer
	if len(prefix) > 0 {
		r.outputPrefix = prefix[0]
	}
	return r
}

// WithErrorOutput
func (r *ScriptBuilder) WithErrorOutput(writer io.Writer) *ScriptBuilder {
	r.errorOutput = writer
	return r
}

// Result runs the script and returns a ScriptResult object.
// If and output writer was provided will automatically write the output
// of the script into the writer.
func (r *ScriptBuilder) Result() *ScriptResult {
	args := r.args
	if r.repoPath != "" {
		args = append([]string{r.repoPath}, args...)
	}
	result := ExecBashScript(path.Join(r.rootFolder, r.script), args...)

	if r.output != nil {
		r.output.Write([]byte(result.OutputData(r.outputPrefix)))
	}
	if r.errorOutput != nil {
		r.errorOutput.Write([]byte(result.Stderr()))
	}
	if result.ExitCode() != 0 {
		logger := logging.FromContext(r.ctx)
		logger.Errorf("result %s returned non-zero code. \nstdout: %s \nstderr: %s", r.script, result.Stdout(), result.Stderr())
	}
	return result
}

// Error runs the script and returns an error if any.
// This method will discard the ScriptResult object
// but the output can be written providing a writer
// using WithOutput method
func (r *ScriptBuilder) Error() error {
	result := r.Result()
	if result.ExitCode() != 0 {
		return errors.New(fmt.Sprintf("the script %s returned an non-zero code %d.\nstdout:%s\nstderr:%s\nerr: %s", r.script, result.ExitCode(), result.Stdout(), result.Stderr(), result.Error()))
	}
	return nil
}
