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

type LabelLevel string

const (
	ModuleLevel    LabelLevel = "module"
	FeatureLevel   LabelLevel = "feature"
	TestCaseLevel  LabelLevel = "testcase"
	TestPointLevel LabelLevel = "testpoint"
)

func NewModuleLabel(moduleName string) *LevelLabel {
	return &LevelLabel{Level: ModuleLevel, Name: moduleName}
}

func NewFeatureLabel(featureName string) *LevelLabel {
	return &LevelLabel{Level: FeatureLevel, Name: featureName}
}

func newTestCaseLabel(testCaseName string) *LevelLabel {
	return &LevelLabel{Level: TestCaseLevel, Name: testCaseName}
}

func newTestPointLabel(testPointName string) *LevelLabel {
	return &LevelLabel{Level: TestPointLevel, Name: testPointName}
}

type LevelLabel struct {
	Level LabelLevel `json:"level" yaml:"level"`
	Name  string     `json:"name" yaml:"name"`
}

func (t LevelLabel) String() string {
	return fmt.Sprintf("%s:%s", t.Level, t.Name)
}

func (t LevelLabel) Labels() Labels {
	return Labels{t.String()}
}
