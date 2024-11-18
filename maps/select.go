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

package maps

// MatchFunc return true if key and value match condition
type MatchFunc func(key, value string) bool

// MutateFunc return mutated key and value
type MutateFunc func(key, value string) (string, string)

// SelectAndMutateMap select elements from a map based on a match function and mutate it's key and value
func SelectAndMutateMap(maps map[string]string, match MatchFunc, mutate MutateFunc) map[string]string {
	result := make(map[string]string)
	for key, value := range maps {
		if match(key, value) {
			key, value = mutate(key, value)
			result[key] = value
		}
	}
	return result
}
