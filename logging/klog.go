/*
Copyright 2024 The AlaudaDevops Authors.

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

package logging

import "strconv"

// KlogLevelKey the key to set klog log level in the configmap
// Example:
//
// data:
//
//	klog.level: "9"
const KlogLevelKey = "klog.level"

// GetKlogLevelFromConfigMapData get klog level from configmap data
func GetKlogLevelFromConfigMapData(data map[string]string) (level string) {
	if data == nil {
		return "0"
	}

	level = data[KlogLevelKey]
	if level == "" {
		return "0"
	}

	_, err := strconv.Atoi(level)
	if err != nil {
		return "0"
	}
	return level
}
