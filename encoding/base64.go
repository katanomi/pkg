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

package encoding

import (
	"encoding/base64"
	"encoding/json"
)

// Base64Encode convert obj to json string, and then encode it with base64.
func Base64Encode(obj interface{}) (string, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(objBytes), nil
}

// Base64Decode decode encodeStr using base64, and then use json to convert to obj.
func Base64Decode(encodeStr string, obj interface{}) error {
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
