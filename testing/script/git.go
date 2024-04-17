/*
Copyright 2023 The Katanomi Authors.

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
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"knative.dev/pkg/logging"
)

//go:embed git
var scripts embed.FS

var e2eTempDir string

const (
	// ScriptInitRepo script path for init repo script
	ScriptInitRepo = "git/init_repo.sh"
	// ScriptCreateBranch script path for create branch script
	ScriptCreateBranch = "git/create_branch.sh"
	// ScriptCreateCommit script path for create commit script
	ScriptCreateCommit = "git/create_commit.sh"
	// ScriptCreateMultiCommit script path for create multi commit script
	ScriptCreateMultiCommit = "git/create_multi_commit.sh"
	// ScriptCreateTag script path for create tag script
	ScriptCreateTag = "git/create_tag.sh"
)

func init() {
	var err error
	e2eTempDir = os.Getenv("E2E_TEMP_DIR")
	if e2eTempDir == "" {
		e2eTempDir, err = os.MkdirTemp("", "e2e")
	} else {
		err = os.MkdirAll(e2eTempDir, 0750)
	}

	if err != nil {
		panic(fmt.Sprintf("failed to create directory of e2e scripts: %s", err))
	}

	RestoreDirectories(scripts, "git", e2eTempDir)
}

// InitRepo init repo, e.g: create default branch, submit some files
// TODO: move to use script builder
func InitRepo(ctx context.Context, repoUrl, username, password string) (repoPath string, err error) {
	result := ExecBashScript(path.Join(e2eTempDir, ScriptInitRepo), repoUrl, username, password)
	if result.ExitCode() != 0 {
		logger := logging.FromContext(ctx)
		logger.Errorf("init repo script returned non-zero code. \nstdout: %s \nstderr: %s", result.Stdout(), result.Stderr())
		return "", errors.New(fmt.Sprintf("failed to init repo, err: %s", result.ExitMessage()))
	}

	return result.OutputData(""), nil
}

// CreateNewBranch create new branch using a bash script and vanila git cli
// returns an error if any
func CreateNewBranch(ctx context.Context, branchName string) (err error) {
	return NewScript(ctx, ScriptCreateBranch, branchName).Error()
}

// CreateNewCommit create new commit using a bash script and vanila git cli
// returns an error if any
func CreateNewCommit(ctx context.Context, branchName, message string) (commitId string, err error) {
	if message == "" {
		message = "e2e commit"
	}
	outputBuffer := &bytes.Buffer{}
	err = NewScript(ctx, ScriptCreateCommit, branchName, message).WithOutput(outputBuffer).Error()
	commitId = outputBuffer.String()
	return
}

// CreateMultiCommit create new commit using a bash script and vanila git cli
// returns an error if any
func CreateMultiCommit(ctx context.Context, branchName, message string, quantity int) (err error) {
	if message == "" {
		message = "e2e commit"
	}
	return NewScript(ctx, ScriptCreateMultiCommit, branchName, message, strconv.Itoa(quantity)).Error()
}

// CreateNewTag create new tag using a bash script and vanila git cli
// returns an error if any
func CreateNewTag(ctx context.Context, branchName, message, tag string) (err error) {
	return NewScript(ctx, ScriptCreateTag, branchName, message, tag).Error()
}
