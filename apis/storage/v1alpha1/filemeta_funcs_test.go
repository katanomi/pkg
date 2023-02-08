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

package v1alpha1

import "testing"

func TestParseFileKey(t *testing.T) {

	tests := []struct {
		name               string
		key                string
		wantPluginName     string
		wantFileObjectName string
	}{
		{
			name:               "empty key",
			key:                "",
			wantPluginName:     "",
			wantFileObjectName: "",
		},
		{
			name:               "invalid key",
			key:                "xxx",
			wantPluginName:     "",
			wantFileObjectName: "",
		},
		{
			name:               "valid key",
			key:                "xxx:/xxx/yyy",
			wantPluginName:     "xxx",
			wantFileObjectName: "/xxx/yyy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPluginName, gotFileObjectName := ParseFileKey(tt.key)
			if gotPluginName != tt.wantPluginName {
				t.Errorf("ParseFileKey() gotPluginName = %v, want %v", gotPluginName, tt.wantPluginName)
			}
			if gotFileObjectName != tt.wantFileObjectName {
				t.Errorf("ParseFileKey() gotFileObjectName = %v, want %v", gotFileObjectName, tt.wantFileObjectName)
			}
		})
	}
}
