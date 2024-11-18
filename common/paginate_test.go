/*
Copyright 2021 The AlaudaDevops Authors.

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

package common

import "testing"

func TestPaginate(t *testing.T) {
	type args struct {
		total   int
		perPage int
		page    int
	}
	tests := []struct {
		name      string
		args      args
		wantBegin int
		wantEnd   int
	}{
		{
			name:      "normal paginate",
			args:      args{total: 88, perPage: 20, page: 3},
			wantBegin: 40,
			wantEnd:   60,
		},
		{
			name:      "input zero paginate",
			args:      args{total: 88, perPage: 0, page: 0},
			wantBegin: 0,
			wantEnd:   20,
		},
		{
			name:      "out total paginate",
			args:      args{total: 88, perPage: 40, page: 3},
			wantBegin: 80,
			wantEnd:   88,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBegin, gotEnd := Paginate(tt.args.total, tt.args.perPage, tt.args.page)
			if gotBegin != tt.wantBegin {
				t.Errorf("Paginate() gotBegin = %v, want %v", gotBegin, tt.wantBegin)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("Paginate() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
