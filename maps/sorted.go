/*
Copyright 2022 The Katanomi Authors.

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

import "sort"

// KeyValue a key/value pair just like a map[string]string
type KeyValue struct {
	// Key of the map
	Key string
	// Value of string map
	Value string
}

// SortedKeyValue returns a list of key values sorted by key
func SortedKeyValue(dict map[string]string) (items []KeyValue) {
	items = make([]KeyValue, 0, len(dict))
	for k, v := range dict {
		items = append(items, KeyValue{Key: k, Value: v})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})
	return
}
