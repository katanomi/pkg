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
