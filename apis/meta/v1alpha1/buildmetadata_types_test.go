package v1alpha1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildMetaData_Serialization(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Status     BuildMetaDataStatus
	}
	type args struct {
		s       string
		replace bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Success", fields: fields{}, args: args{
			s:       `{\"kind\":\"BuildRun\",\"apiVersion\":\"builds.katanomi.dev/v1alpha1\",\"metadata\":{\"name\":\"build-223f80-2xgnq\",\"namespace\":\"kychen\",\"creationTimestamp\":null,\"labels\":{\"builds.katanomi.dev/build\":\"build-223f80\"},\"annotations\":{\"integrations.katanomi.dev/clusterName\":\"global\",\"integrations.katanomi.dev/integration\":\"gitlab-kychen\",\"katanomi.dev/createdBy\":\"{\\\"user\\\":{\\\"kind\\\":\\\"User\\\",\\\"name\\\":\\\"admin@cpaas.io\\\"}}\"}},\"status\":{\"git\":{\"branch\":{\"name\":\"build\"},\"lastCommit\":{\"id\":\"cad87afa9016c50489cbd2f48b3199879224ad7c\"}},\"triggeredBy\":{\"createdBy\":{\"user\":{\"kind\":\"User\",\"name\":\"admin@cpaas.io\"}}}}}`,
			replace: false,
		},
			wantErr: true,
		},
		{name: "Success nomal", fields: fields{}, args: args{
			s:       `{"kind":"BuildRun","apiVersion":"builds.katanomi.dev/v1alpha1","metadata":{"name":"build-223f80-2xgnq","namespace":"kychen","creationTimestamp":null,"labels":{"builds.katanomi.dev/build":"build-223f80"},"annotations":{"integrations.katanomi.dev/clusterName":"global","integrations.katanomi.dev/integration":"gitlab-kychen","katanomi.dev/createdBy":"{\"user\":{\"kind\":\"User\",\"name\":\"admin@cpaas.io\"}}"}},"status":{"git":{"branch":{"name":"build"},"lastCommit":{"id":"cad87afa9016c50489cbd2f48b3199879224ad7c"}},"triggeredBy":{"createdBy":{"user":{"kind":"User","name":"admin@cpaas.io"}}}}}`,
			replace: false,
		},
			wantErr: false,
		},
		{name: "Success replace", fields: fields{}, args: args{
			s:       `{"kind":"BuildRun","apiVersion":"builds.katanomi.dev/v1alpha1","metadata":{"name":"build-223f80-6x7c5","namespace":"kychen","creationTimestamp":null,"labels":{"builds.katanomi.dev/build":"build-223f80"},"annotations":{"integrations.katanomi.dev/clusterName":"global","integrations.katanomi.dev/integration":"gitlab-kychen","katanomi.dev/createdBy":"{\\"user\\":{\\"kind\\":\\"User\\",\\"name\\":\\"admin@cpaas.io\\"}}"}},"status":{"git":{"branch":{"name":"build"},"lastCommit":{"id":"cad87afa9016c50489cbd2f48b3199879224ad7c"}},"triggeredBy":{"createdBy":{"user":{"kind":"User","name":"admin@cpaas.io"}}}}}`,
			replace: true,
		},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &BuildMetaData{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Status:     tt.fields.Status,
			}
			if err := p.Serialization(tt.args.s, tt.args.replace); (err != nil) != tt.wantErr {
				t.Errorf("BuildMetaData.Serialization() \nerror = %v, \n wantErr %v", err, tt.wantErr)
			}
		})
	}
}
