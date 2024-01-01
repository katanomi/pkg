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

package conformance

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
)

func NewNode(level LabelLevel, name string) *Node {
	return &Node{
		LevelLabel: &LevelLabel{
			Level: level,
			Name:  name,
		},
	}
}

// Node is the basic unit of the test case set
type Node struct {
	*LevelLabel `json:",inline" yaml:",inline"`
	Description string `json:"description" yaml:"description"`

	caseRegister func()

	ParentNode *Node   `json:"parentNode,omitempty" yaml:"parentNode,omitempty"`
	SubNodes   []*Node `json:"subNodes" yaml:"subNodes"`
}

// RegisterTestCase iterate over the node tree, register all the test case to ginkgo
func (n *Node) RegisterTestCase() {
	n.registerTestCase(n)
}

func (n *Node) registerTestCase(node *Node) {
	Describe(fmt.Sprintf("test for %s %s", node.Name, node.Level), node.IdentifyLabels(), func() {
		if node.caseRegister != nil {
			node.caseRegister()
		}

		if len(node.SubNodes) > 0 {
			for _, subNode := range node.SubNodes {
				n.registerTestCase(subNode)
			}
		}
	})
}

// Clone deep clone the node and sub nodes
func (n *Node) Clone() *Node {
	return n.clone(n)
}

func (n *Node) clone(original *Node) *Node {
	clone := *original

	if len(original.SubNodes) > 0 {
		clone.SubNodes = make([]*Node, len(original.SubNodes))
		for i, subNode := range original.SubNodes {
			clone.SubNodes[i] = n.clone(subNode)
			clone.SubNodes[i].ParentNode = &clone
		}
	}

	return &clone
}

// FullPathLabels return all the full path to the leaf nodeS
func (n *Node) FullPathLabels() Labels {
	var paths []string
	var path []string
	n.traverse(n, path, &paths)
	return paths
}

func (n *Node) traverse(node *Node, path []string, paths *[]string) {
	path = append(path, node.String())

	if len(node.SubNodes) == 0 {
		*paths = append(*paths, strings.Join(path, "#"))
	} else {
		for _, subNode := range node.SubNodes {
			n.traverse(subNode, path, paths)
		}
	}

	path = path[:len(path)-1]
}

// Equal check if two nodes are equal
func (n *Node) Equal(other *Node) bool {
	return n.Level == other.Level && n.Name == other.Name
}

// LinkParentNode link current node to a parent node
func (n *Node) LinkParentNode(parentNode *Node) {
	n.ParentNode = parentNode
	parentNode.SubNodes = append(parentNode.SubNodes, n)
}

// IdentifyLabels return the unique labels of the current node
func (n *Node) IdentifyLabels() Labels {
	return Labels{strings.Join(n.Labels(), "#")}
}

// Labels return all the labels contains all parent node labels
func (n *Node) Labels() Labels {
	labels := n.LevelLabel.Labels()
	if n.ParentNode != nil {
		labels = append(n.ParentNode.Labels(), labels...)
	}
	return labels
}
