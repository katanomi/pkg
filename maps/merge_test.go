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

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMergeMap(t *testing.T) {

	table := map[string]struct {
		Left   map[string]string
		Right  map[string]string
		Result map[string]string
	}{
		"multiple keys and values": {
			Left: map[string]string{
				"a": "b",
			},
			Right: map[string]string{
				"b": "c",
				"a": "d",
			},
			Result: map[string]string{
				"b": "c",
				"a": "d",
			},
		},
		"nil left": {
			Left: nil,
			Right: map[string]string{
				"b": "c",
				"a": "d",
			},
			Result: map[string]string{
				"b": "c",
				"a": "d",
			},
		},
		"nil right": {
			Left: map[string]string{
				"a": "b",
			},
			Right: nil,
			Result: map[string]string{
				"a": "b",
			},
		},
		"both nil": {
			Left:   nil,
			Right:  nil,
			Result: map[string]string{},
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := MergeMap(test.Left, test.Right)
			g.Expect(result).To(Equal(test.Result))
		})
	}
}

func TestMergeMapIfNotExists(t *testing.T) {

	table := map[string]struct {
		Left   map[string]string
		Right  map[string]string
		Result map[string]string
	}{
		"multiple keys and values": {
			Left: map[string]string{
				"a": "b",
				"d": "",
			},
			Right: map[string]string{
				"a": "d",
				"b": "c",
				"c": "c1",
				"d": "d1",
			},
			Result: map[string]string{
				"a": "b",
				"b": "c",
				"c": "c1",
				"d": "",
			},
		},
		"nil left": {
			Left: nil,
			Right: map[string]string{
				"b": "c",
				"a": "d",
			},
			Result: map[string]string{
				"b": "c",
				"a": "d",
			},
		},
		"nil right": {
			Left: map[string]string{
				"a": "b",
			},
			Right: nil,
			Result: map[string]string{
				"a": "b",
			},
		},
		"both nil": {
			Left:   nil,
			Right:  nil,
			Result: nil,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := MergeMapIfNotExists(test.Left, test.Right)
			g.Expect(result).To(Equal(test.Result))
		})
	}
}

func TestMergeMapSlice(t *testing.T) {

	table := map[string]struct {
		Left   map[string][]string
		Right  map[string][]string
		Result map[string][]string
	}{
		"multiple keys and values": {
			Left: map[string][]string{
				"a": {"b", "c"},
			},
			Right: map[string][]string{
				"b": {"c", "d"},
				"a": {"d", "e"},
			},
			Result: map[string][]string{
				"b": {"c", "d"},
				"a": {"d", "e"},
			},
		},
		"nil left": {
			Left: nil,
			Right: map[string][]string{
				"b": {"c", "d"},
				"a": {"d", "e"},
			},
			Result: map[string][]string{
				"b": {"c", "d"},
				"a": {"d", "e"},
			},
		},
		"nil right": {
			Left: map[string][]string{
				"x": {"z", "q"},
				"y": {"k", "w"},
			},
			Right: nil,
			Result: map[string][]string{
				"x": {"z", "q"},
				"y": {"k", "w"},
			},
		},
		"both nil": {
			Left:   nil,
			Right:  nil,
			Result: map[string][]string{},
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := MergeMapSlice(test.Left, test.Right)
			g.Expect(result).To(Equal(test.Result))
		})
	}
}

func TestMergeMapMap(t *testing.T) {

	table := map[string]struct {
		Left   map[string]map[string]string
		Right  map[string]map[string]string
		Result map[string]map[string]string
	}{
		"multiple keys and values": {
			Left: map[string]map[string]string{
				"a": {"b": "c"},
			},
			Right: map[string]map[string]string{
				"b": {"c": "d"},
				// b key in a will be replaced
				"a": {"d": "e", "b": "y"},
			},
			Result: map[string]map[string]string{
				"b": {"c": "d"},
				"a": {"d": "e", "b": "y"},
			},
		},
		"left nil": {
			Left: nil,
			Right: map[string]map[string]string{
				"b": {"c": "d"},
				"a": {"d": "e", "b": "y"},
			},
			Result: map[string]map[string]string{
				"b": {"c": "d"},
				"a": {"d": "e", "b": "y"},
			},
		},
		"right nil": {
			Left: map[string]map[string]string{
				"a": {"b": "c"},
			},
			Right: nil,
			Result: map[string]map[string]string{
				"a": {"b": "c"},
			},
		},
		"both nil": {
			Left:   nil,
			Right:  nil,
			Result: map[string]map[string]string{},
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := MergeMapMap(test.Left, test.Right)
			g.Expect(result).To(Equal(test.Result))
		})
	}
}
