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

// Package common useful functionality for encoding and decoding objects
package common

import (
	"encoding/base64"
	"encoding/json"
)

// ToJSONBase64 convert obj to json string, and then encode it with base64.
func ToJSONBase64(obj interface{}) (string, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(objBytes), nil
}

// FromJSONBase64 decode encodeStr using base64, and then use json to convert to obj.
func FromJSONBase64(encodeStr string, obj interface{}) error {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodeBytes, obj)
	if err != nil {
		return err
	}

	return nil
}
