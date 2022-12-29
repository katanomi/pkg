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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// ExitCodeOption adds a generic option to store different cli paths
type ExitCodeOption struct {
	// ExitCodePath direct path for exitCodefile
	ExitCodePath string
	// FlagName to be used to store
	FlagName string
}

// AddFlags adds flags for option
func (p *ExitCodeOption) AddFlags(flags *pflag.FlagSet) {
	if p.FlagName == "" {
		p.FlagName = "exit-code-path"
	}
	flags.StringVar(&p.ExitCodePath, p.FlagName, "", "specify the exit status code to save the file.")
}

// Succeed judging whether the execution status is successful or not according to the exit file.
func (m *ExitCodeOption) Succeed() (bool, string, error) {
	if m.ExitCodePath == "" {
		return false, "", fmt.Errorf("exit code file not set")
	}

	data, err := os.ReadFile(m.ExitCodePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to read exit code file[%s], error: %s", m.ExitCodePath, err.Error())
	}

	code := string(data)
	if code != "" && strings.TrimSpace(code) != "0" {
		return false, code, nil
	}
	return true, code, nil
}
