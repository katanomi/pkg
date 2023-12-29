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

package client

import "github.com/katanomi/pkg/apis/meta/v1alpha1"

var _ PluginDisplayColumns = &DisplayColumnsManager{}

// DisplayColumnsManager provide default displayColumn management.
// When the default implementation is not satisfied, you need to implement it yourself.
type DisplayColumnsManager struct {
	displayColumns map[string]v1alpha1.DisplayColumns
}

// SetDisplayColumns set display columns
func (d *DisplayColumnsManager) SetDisplayColumns(key string, values ...v1alpha1.DisplayColumn) {
	if d.displayColumns == nil {
		d.displayColumns = make(map[string]v1alpha1.DisplayColumns)
	}
	d.displayColumns[key] = values
}

// GetDisplayColumns get display columns
func (d *DisplayColumnsManager) GetDisplayColumns() map[string]v1alpha1.DisplayColumns {
	if d.displayColumns == nil {
		return make(map[string]v1alpha1.DisplayColumns)
	}
	return d.displayColumns
}
