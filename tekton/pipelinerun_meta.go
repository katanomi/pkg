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

package tekton

import (
	"encoding/base64"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/tools/reference"
	"sigs.k8s.io/yaml"
)

// PipelineRunMeta defines fields injecting to pipelinerun
type PipelineRunMeta struct {
	// RunRefs are owner's of pipelinerun from top to bottom
	RunRefs []*v1.ObjectReference `json:"runRefs,omitempty"`
}

func (m *PipelineRunMeta) Encode() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// InjectPipelineRunMeta injects encoded meta.json into pipelinerun's annotation
func InjectPipelineRunMeta(pr *v1beta1.PipelineRun, meta PipelineRunMeta) error {
	if pr == nil {
		return reference.ErrNilObject
	}

	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	dataEncoded := base64.StdEncoding.EncodeToString(data)

	annotations := pr.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string, 0)
	}
	annotations[PodAnnotationMetaJson] = dataEncoded
	pr.SetAnnotations(annotations)
	return nil
}

// UnmarshalMeta unmarshal pipelinerun meta
func UnmarshalMeta(metaJsonString string) (*PipelineRunMeta, error) {
	decodeString, err := base64.StdEncoding.DecodeString(metaJsonString)
	if err != nil {
		return nil, err
	}
	decodedMeta := &PipelineRunMeta{}
	err = yaml.Unmarshal(decodeString, &decodedMeta)
	if err != nil {
		return nil, err
	}
	return decodedMeta, nil
}
