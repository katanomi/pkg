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

package encoding

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var sliceRegx = regexp.MustCompile(`\[(\d+)\]$`)

type NodeType string

const (
	NodeTypeRoot   NodeType = "root"
	NodeTypeArray  NodeType = "array"
	NodeTypeObject NodeType = "object"
	NodeTypeValue  NodeType = "value"
)

func NewRootNode() *Node {
	return &Node{
		Type:   NodeTypeRoot,
		Childs: map[string]*Node{},
	}
}

type Node struct {
	Name   string
	Value  string
	Type   NodeType
	Childs map[string]*Node
}

func (n *Node) sliceToJson() string {
	s := &strings.Builder{}
	s.WriteByte('[')

	// sort by key
	keys := make([]int, 0, len(n.Childs))
	for key := range n.Childs {
		index, _ := strconv.Atoi(key)
		keys = append(keys, index)
	}
	sort.Ints(keys)

	index := 0
	for _, key := range keys {
		child := n.Childs[strconv.Itoa(key)]
		s.WriteString(child.Json())
		if index++; index != len(n.Childs) {
			s.WriteByte(',')
		}
	}
	s.WriteByte(']')
	return s.String()
}

func (n *Node) objectToJson() string {
	s := &strings.Builder{}
	s.WriteByte('{')
	index := 0
	for key, child := range n.Childs {
		s.WriteString(fmt.Sprintf(`"%s":%s`, key, child.Json()))
		index++
		if index != len(n.Childs) {
			s.WriteByte(',')
		}
	}
	s.WriteByte('}')
	return s.String()
}

func (n *Node) Json() string {
	switch n.Type {
	case NodeTypeArray:
		return n.sliceToJson()
	case NodeTypeRoot, NodeTypeObject:
		return n.objectToJson()
	case NodeTypeValue:
		return fmt.Sprintf("\"%s\"", n.Value)
	default:
		return ""
	}
}

func (n *Node) Set(path string, value string) {
	assignNodeValue(strings.Split(path, "."), value, n)
}

func (n *Node) Get(paths []string) *Node {
	if n == nil {
		return nil
	}
	if len(paths) == 0 {
		return n
	}
	path := paths[0]
	// slice
	if list := sliceRegx.FindStringSubmatch(path); len(list) > 0 {
		path = sliceRegx.ReplaceAllString(path, "")
		paths = append([]string{path, list[1]}, paths[1:]...)
		return n.Get(paths)
	}

	subNode := n.Childs[path]
	if subNode == nil {
		return nil
	}
	return subNode.Get(paths[1:])
}

func assignNodeValue(paths []string, value string, node *Node) {
	if len(paths) == 0 || node == nil {
		return
	}

	// slice
	// eg: list[0]
	path := paths[0]
	if list := sliceRegx.FindStringSubmatch(path); len(list) > 0 {
		path = sliceRegx.ReplaceAllString(path, "")
		if path == "" {
			// invalid or not support, ignore
			return
		}
		subNode := node.Get([]string{path})
		if subNode == nil {
			subNode = &Node{
				Name:   path,
				Type:   NodeTypeArray,
				Childs: make(map[string]*Node),
			}
			node.Childs[subNode.Name] = subNode
		}
		assignNodeValue(append([]string{list[1]}, paths[1:]...), value, subNode)
		return
	}

	// value
	if len(paths) == 1 {
		node.Childs[path] = &Node{
			Name:  paths[0],
			Type:  NodeTypeValue,
			Value: value,
		}
		return
	}

	// object
	subNode := node.Get([]string{path})
	if subNode == nil {
		subNode = &Node{
			Name:   path,
			Type:   NodeTypeObject,
			Childs: make(map[string]*Node),
		}
		node.Childs[subNode.Name] = subNode
	}
	assignNodeValue(paths[1:], value, subNode)
}
