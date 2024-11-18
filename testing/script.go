/*
Copyright 2023 The AlaudaDevops Authors.

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

package testing

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// DefaultScriptOutputPrefix the prefix of script output
// Deprecated: move to testing/script DefaultScriptOutputPrefix
const DefaultScriptOutputPrefix = "##output##"

// RestoreDirectories restore directories from embed.FS to targetDir
// Deprecated: move to testing/script RestoreDirectories
func RestoreDirectories(fs embed.FS, dirName string, targetDir string) {
	entries, err := fs.ReadDir(dirName)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			RestoreDirectories(fs, path.Join(dirName, entry.Name()), targetDir)
		} else {
			filePath := path.Join(dirName, entry.Name())
			targetFilePath := path.Join(targetDir, filePath)
			fileContent, _ := fs.ReadFile(filePath)
			os.MkdirAll(path.Dir(targetFilePath), 0750)
			err = os.WriteFile(targetFilePath, fileContent, 0750)
			if err != nil {
				panic(fmt.Sprintf("failed to save %s to %s, err: %s", filePath, targetFilePath, err))
			}
		}
	}
}

// ExecBashScript exec bash script with params
// Deprecated: move to testing/script ExecBashScript
func ExecBashScript(script string, params ...string) *ScriptResult {
	return ExecScript("bash", append([]string{script}, params...)...)
}

// ExecScript exec script with params
// Deprecated: move to testing/script ExecScript
func ExecScript(name string, arg ...string) *ScriptResult {
	c := exec.Command(name, arg...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()

	result := &ScriptResult{
		stdout: stdout,
		stderr: stderr,
	}
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.exitError = exitError
		}
	}

	return result
}

// ScriptResult the result of script execution
// Deprecated: move to testing/script ScriptResult
type ScriptResult struct {
	stdout    bytes.Buffer
	stderr    bytes.Buffer
	exitError *exec.ExitError

	OutputDataPrefix string
}

// Error get the error of a script execution
func (p *ScriptResult) Error() error {
	return p.exitError
}

// Stdout get the output of a script execution
func (p *ScriptResult) Stdout() string {
	return p.stdout.String()
}

// Stderr get the error output of a script execution
func (p *ScriptResult) Stderr() string {
	return p.stderr.String()
}

// ExitCode get the exit code of a script execution
func (p *ScriptResult) ExitCode() int {
	if p.exitError == nil {
		return 0
	}
	return p.exitError.ExitCode()
}

// ExitMessage get the exit message of a script execution
func (p *ScriptResult) ExitMessage() string {
	if p.exitError == nil {
		return ""
	}
	return p.exitError.String()
}

// OutputData get the output of a script execution,
// it can be a structured data, e.g: json string
func (p *ScriptResult) OutputData(dataPrefix string) string {
	if dataPrefix == "" {
		dataPrefix = DefaultScriptOutputPrefix
	}
	parts := strings.Split(p.Stdout(), "\n")
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasPrefix(parts[i], dataPrefix) {
			return strings.TrimPrefix(parts[i], dataPrefix)
		}
	}
	return ""
}
