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
	"strings"

	"github.com/katanomi/pkg/testing/framework/base"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
)

var globalModuleCases []*moduleCase
var Configure = base.ConfigureFunc(func(suiteConfig *types.SuiteConfig, reporterConfig *types.ReporterConfig) {
	suiteConfig.LabelFilter = strings.Join(GlobalModuleCaseLabels(), ",")
})

func pushModuleCase(m *moduleCase) {
	globalModuleCases = append(globalModuleCases, m)
}

// GlobalModuleCaseLabels get all labels of conformance testcases
func GlobalModuleCaseLabels() Labels {
	var labels Labels
	for _, item := range globalModuleCases {
		labels = append(labels, item.node.FullPathLabels()...)
	}

	config, _ := GinkgoConfiguration()
	labelFilter, err := types.ParseLabelFilter(config.LabelFilter)
	if err != nil {
		return labels
	}

	var focusLabels []string
	for _, item := range labels {
		if labelFilter(Labels{item}) {
			focusLabels = append(focusLabels, item)
		}
	}

	return focusLabels
}
