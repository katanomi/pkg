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

package encoding

import (
	"reflect"
	"testing"
)

func TestToJSONBase64_FromJSONBase64(t *testing.T) {
	type TestStruct struct {
		URL  string `json:"url,omitempty"`
		Name string `json:"name,omitempty"`
	}
	type args struct {
		obj *TestStruct
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ToJSONBase64 OK",
			args: args{
				obj: &TestStruct{URL: "https://book.kubebuilder.io/", Name: "!2@3\\(0&"},
			},
			want:    "eyJ1cmwiOiJodHRwczovL2Jvb2sua3ViZWJ1aWxkZXIuaW8vIiwibmFtZSI6IiEyQDNcXCgwXHUwMDI2In0=",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Encode(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSONBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJSONBase64() = %v, want %v", got, tt.want)
			}

			out := &TestStruct{}
			if err := Base64Decode(got, out); (err != nil) != tt.wantErr {
				t.Errorf("FromJSONBase64() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(out, tt.args.obj) {
				t.Errorf("FromJSONBase64() = %v, want %v", out, tt.args.obj)
			}
		})
	}
}
