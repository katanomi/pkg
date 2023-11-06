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
	"errors"
	"fmt"
	"reflect"

	"github.com/katanomi/pkg/plugin/client"
	. "github.com/katanomi/pkg/testing/framework/base"
)

// TestablePlugin describe the interface of a testable plugin
type TestablePlugin interface {
	client.PluginRegister

	// GetTestOrgProject get the project name with org type for testing
	GetTestOrgProject() string

	// GetTestUserProject get the project name with user type for testing
	GetTestUserProject() string
}

// PluginTestCase is a test case
type PluginTestCase = func()

// TestSpec test spec which describe the detail of a test case
type TestSpec func(testCtx context.Context, instance TestablePlugin)

// TestCaseGenerator generate a test case
type TestCaseGenerator func(testCtx context.Context, ins TestablePlugin) PluginTestCase

type PluginImplementCondition struct {
	Interface interface{}
}

func (p PluginImplementCondition) Condition(testCtx *TestContext) error {
	testPlugin := GitPluginFromCtx(testCtx.Context)
	ifaceType := reflect.TypeOf(p.Interface).Elem()
	if !reflect.TypeOf(testPlugin).Implements(ifaceType) {
		return errors.New(fmt.Sprintf("plugin does not implement %s interface", ifaceType.Name()))
	}
	return nil
}
