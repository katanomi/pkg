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

package types

import (
	"context"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/test_module.go github.com/katanomi/pkg/plugin/types TestModuleLister

// TestModuleLister list a test module
type TestModuleLister interface {
	Interface
	ListTestModules(ctx context.Context, params metav1alpha1.TestProjectOptions, options metav1alpha1.ListOptions) (*metav1alpha1.TestModuleList, error)
}
