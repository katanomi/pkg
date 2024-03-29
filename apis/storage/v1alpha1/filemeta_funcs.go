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
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// ParseFileKey returns pluginName and fileObjectName
func ParseFileKey(key string) (pluginName, fileObjectName string) {
	keySplit := strings.Split(key, ":")
	if len(keySplit) == 2 {
		return keySplit[0], keySplit[1]
	}
	return "", ""
}

// Encode returns json base64 encoded string
func (in *FileMeta) Encode() string {
	marshaledMeta, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(marshaledMeta)
}

// LastModified return lastModifiedTime from storage server or zero time if no creation info
func (in *FileMeta) LastModified() time.Time {
	if lastModifiedAnno, ok := in.Annotations[LastModifiedAnnotation]; ok {
		parsedTime, _ := time.Parse(time.RFC3339, lastModifiedAnno)
		return parsedTime
	}
	return time.Time{}
}

// Key returns file key
func (in *FileMeta) Key() string {
	return in.Spec.Key
}

// DecodeAsFileMeta decodes encoded string to a pointer FileMeta
func DecodeAsFileMeta(encoded string) (*FileMeta, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	fileMeta := &FileMeta{}
	err = json.Unmarshal(decodedBytes, fileMeta)
	if err != nil {
		return nil, err
	}
	return fileMeta, nil
}
