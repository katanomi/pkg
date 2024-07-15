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
	"reflect"
	"testing"
)

func TestConvertMedataSelectorToConditions(t *testing.T) {
	type args struct {
		metadataSelectorStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []Condition
		wantErr bool
	}{
		{
			name: "equal and like operator",
			args: args{
				metadataSelectorStr: "environment = production,topic=like--consumer",
			},
			want: []Condition{
				{
					Key:      MetadataKey("environment"),
					Operator: ConditionOperatorEqual,
					Value:    "production",
				},
				{
					Key:      MetadataKey("topic"),
					Operator: ConditionOperatorLike,
					Value:    "consumer",
				},
			},
		},
		{
			name: "in operator",
			args: args{
				metadataSelectorStr: "project in (demo1,demo2)",
			},
			want: []Condition{
				{
					Key:      "",
					Operator: ConditionOperatorOr,
					Value: []Condition{
						{
							Key:      MetadataKey("project"),
							Operator: ConditionOperatorEqual,
							Value:    "demo1",
						},
						{
							Key:      MetadataKey("project"),
							Operator: ConditionOperatorEqual,
							Value:    "demo2",
						},
					},
				},
			},
		},
		{
			name: "exist operator",
			args: args{
				metadataSelectorStr: "partition",
			},
			want: []Condition{
				{
					Key:      MetadataKey("partition"),
					Operator: ConditionOperatorExist,
					Value:    "",
				},
			},
		},
		{
			name: "not support operator",
			args: args{
				metadataSelectorStr: "!partition",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMedataSelectorToConditions(tt.args.metadataSelectorStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertMedataSelectorToConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertMedataSelectorToConditions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
