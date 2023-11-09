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

package gitplugin

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/katanomi/pkg/testing"
)

//go:embed scripts
var scripts embed.FS

var e2eTempDir string

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

	testing.RestoreDirectories(scripts, "scripts", e2eTempDir)
}

// InitRepo init repo, e.g: create default branch, submit some files
func InitRepo(ctx context.Context, repoUrl, username, password string) (repoPath string, err error) {
	result := testing.ExecBashScript(path.Join(e2eTempDir, "scripts/init_repo.sh"), repoUrl, username, password)
	if result.ExitCode() != 0 {
		return "", errors.New(fmt.Sprintf("failed to init repo, err: %s", result.ExitMessage()))
	}
	return result.OutputData(""), nil
}

func CreateByScript(ctx context.Context, scriptName string, params ...string) (err error) {
	localRepoPath := LocalRepoPathFromCtx(ctx)
	if localRepoPath == nil {
		return errors.New("no local repo path found")
	}

	_params := append([]string{*localRepoPath}, params...)
	return testing.ExecBashScript(path.Join(e2eTempDir, scriptName), _params...).Error()
}

// CreateNewBranch create new branch
func CreateNewBranch(ctx context.Context, branchName string) (err error) {
	return CreateByScript(ctx, "scripts/create_branch.sh", branchName)
}

// CreateNewCommit create new commit
func CreateNewCommit(ctx context.Context, branchName, message string) (err error) {
	if message == "" {
		message = "e2e commit"
	}
	return CreateByScript(ctx, "scripts/create_commit.sh", branchName, message)
}

// CreatNewCommit create new commit
func CreateMuiltCommit(ctx context.Context, branchName, message string, quantity int) (err error) {
	if message == "" {
		message = "e2e commit"
	}
	return CreateByScript(ctx, "scripts/create_mulit_commit.sh", branchName, message, strconv.Itoa(quantity))
}

// CreateNewTag create new tag
func CreateNewTag(ctx context.Context, branchName, message, tag string) (err error) {
	return CreateByScript(ctx, "scripts/create_tag.sh", branchName, message, tag)
}
