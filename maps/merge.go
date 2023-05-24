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

// Package maps contains methods to operate maps of all kinds
package maps

// MergeMap merges the right map into left map overwritting any matching keys
func MergeMap(left, right map[string]string) map[string]string {
	if left == nil {
		left = map[string]string{}
	}
	for k, v := range right {
		left[k] = v
	}
	return left
}

// MergeMapIfNotExists merges the right map into left map if right key is not exists in left
func MergeMapIfNotExists(left, right map[string]string) map[string]string {

	if right == nil {
		return left
	}

	for k, v := range right {
		if left == nil {
			left = map[string]string{}
		}

		if _, ok := left[k]; !ok {
			left[k] = v
		}
	}
	return left
}

// MergeMapSlice merges the right map into left map overwritting any matching keys
func MergeMapSlice(left, right map[string][]string) map[string][]string {
	if left == nil {
		left = map[string][]string{}
	}
	for k, v := range right {
		left[k] = v
	}
	return left
}

// MergeMapMap merges the right map into left map overwritting any matching keys
func MergeMapMap(left, right map[string]map[string]string) map[string]map[string]string {
	if left == nil {
		left = map[string]map[string]string{}
	}
	for k, v := range right {
		t := map[string]string{}
		t = MergeMap(t, left[k])
		t = MergeMap(t, v)
		left[k] = t
	}
	return left
}
