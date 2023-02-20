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

import (
	"testing"

	"github.com/katanomi/pkg/pointer"
	"github.com/onsi/gomega"
)

func TestHistoryLimits_IsInvalid(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := map[string]struct {
		limit *HistoryLimits
		want  bool
	}{
		"HistoryLimits is nil": {
			limit: nil,
			want:  false,
		},
		"HistoryLimits.Count is nil": {
			limit: &HistoryLimits{},
			want:  false,
		},
		"HistoryLimits.Count is zero": {
			limit: &HistoryLimits{Count: pointer.Int(0)},
			want:  false,
		},
		"HistoryLimits.Count is -1": {
			limit: &HistoryLimits{Count: pointer.Int(-1)},
			want:  true,
		},
		"HistoryLimits.Count is 1": {
			limit: &HistoryLimits{Count: pointer.Int(1)},
			want:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.limit.IsInvalid()).To(gomega.Equal(tt.want))
		})
	}
}

func TestHistoryLimits_IsNotSet(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := map[string]struct {
		limit *HistoryLimits
		want  bool
	}{
		"HistoryLimits is nil": {
			limit: nil,
			want:  true,
		},
		"HistoryLimits.Count is nil": {
			limit: &HistoryLimits{},
			want:  true,
		},
		"HistoryLimits.Count is zero": {
			limit: &HistoryLimits{Count: pointer.Int(0)},
			want:  false,
		},
		"HistoryLimits.Count is -1": {
			limit: &HistoryLimits{Count: pointer.Int(-1)},
			want:  false,
		},
		"HistoryLimits.Count is 1": {
			limit: &HistoryLimits{Count: pointer.Int(1)},
			want:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.limit.IsNotSet()).To(gomega.Equal(tt.want))
		})
	}
}
