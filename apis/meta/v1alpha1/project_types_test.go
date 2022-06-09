/*
Copyright 2021 The Katanomi Authors.

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
	"reflect"
	"testing"

	authv1 "k8s.io/api/authorization/v1"
)

func TestProjectResourceAttributes(t *testing.T) {
	type args struct {
		verb string
	}
	tests := []struct {
		name string
		args args
		want authv1.ResourceAttributes
	}{
		{name: "ProjectResourceAttributes", args: args{verb: "test"}, want: authv1.ResourceAttributes{Group: "meta.katanomi.dev", Version: "v1alpha1", Resource: "projects", Verb: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProjectResourceAttributes(tt.args.verb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProjectResourceAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
