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

package conformance

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
)

// LabelLevel is the level of the label
type LabelLevel string

const (
	ModuleLevel    LabelLevel = "module"
	FunctionLevel  LabelLevel = "function"
	FeatureLevel   LabelLevel = "feature"
	TestCaseLevel  LabelLevel = "testcase"
	TestPointLevel LabelLevel = "testpoint"
)

func NewModuleLabel(name string) *LevelLabel {
	return &LevelLabel{Level: ModuleLevel, Name: name}
}

func NewFunctionLabel(name string) *LevelLabel {
	return &LevelLabel{Level: FunctionLevel, Name: name}
}

func NewFeatureLabel(name string) *LevelLabel {
	return &LevelLabel{Level: FeatureLevel, Name: name}
}

func newTestCaseLabel(name string) *LevelLabel {
	return &LevelLabel{Level: TestCaseLevel, Name: name}
}

func newTestPointLabel(name string) *LevelLabel {
	return &LevelLabel{Level: TestPointLevel, Name: name}
}

// LevelLabel describe the label of the test node contains level and name
type LevelLabel struct {
	Level LabelLevel `json:"level" yaml:"level"`
	Name  string     `json:"name" yaml:"name"`
}

// String return the string format of the label
func (t LevelLabel) String() string {
	return fmt.Sprintf("%s:%s", t.Level, t.Name)
}

// Labels return the ginkgo-labels format of the label
func (t LevelLabel) Labels() Labels {
	return Labels{t.String()}
}
