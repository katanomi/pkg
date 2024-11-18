/*
Copyright 2024 The AlaudaDevops Authors.

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

package logging

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetKlogLevelFromConfigMapData(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		name      string
		data      map[string]string
		wantLevel string
	}{
		{
			name:      "nil config",
			data:      nil,
			wantLevel: "0",
		},
		{
			name:      "empty config",
			data:      map[string]string{},
			wantLevel: "0",
		},
		{
			name: "not exist the key klog.level",
			data: map[string]string{
				"a": "b",
			},
			wantLevel: "0",
		},
		{
			name: "exist the key klog.level and the value is correct number",
			data: map[string]string{
				KlogLevelKey: "3",
			},
			wantLevel: "3",
		},
		{
			name: "exist the key klog.level and the value is not a number",
			data: map[string]string{
				KlogLevelKey: "a3",
			},
			wantLevel: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(GetKlogLevelFromConfigMapData(tt.data)).To(Equal(tt.wantLevel))
		})
	}
}
