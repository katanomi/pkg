/*
Copyright 2024 The Katanomi Authors.

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
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test for node FullPathLabels", func() {
	var (
		node   *Node
		labels Labels
	)

	BeforeEach(func() {
		node = NewNode(ModuleLevel, "top")
		labels = Labels{}
	})

	JustBeforeEach(func() {
		labels = node.FullPathLabels()
	})

	When("single node", func() {
		It("returned labels should be correct", func() {
			Expect(labels).To(Equal(Labels{"module:top"}))
		})
	})

	When("nested node", func() {
		BeforeEach(func() {
			node1 := NewNode(FeatureLevel, "node1")
			node2 := NewNode(FeatureLevel, "node2")
			node1.LinkParentNode(node)
			node2.LinkParentNode(node1)
		})
		It("returned labels should be correct", func() {
			Expect(labels).To(Equal(Labels{"module:top#feature:node1#feature:node2"}))
		})
	})
})

var _ = Describe("test for node clone", func() {
	var (
		node       *Node
		clonedNode *Node
	)

	BeforeEach(func() {
		node = NewNode(ModuleLevel, "top")
		clonedNode = nil
	})

	compareNode := func(a, b *Node) {
		GinkgoHelper()
		Expect(a.Name).To(Equal(b.Name))
		Expect(a.Level).To(Equal(b.Level))
		Expect(a.ParentNode).To(Equal(b.ParentNode))
		// compare pointers
		Expect(a == b).NotTo(BeTrue())
	}

	checkTopNode := func(node *Node) {
		GinkgoHelper()
		Expect(node).NotTo(BeNil())
		Expect(node.Level).To(Equal(ModuleLevel))
		Expect(node.Name).To(Equal("top"))
	}

	JustBeforeEach(func() {
		clonedNode = node.Clone()
	})

	When("without subnodes", func() {
		It("should be cloned successfully", func() {
			checkTopNode(clonedNode)
			compareNode(clonedNode, node)
			Expect(clonedNode.SubNodes).To(BeNil())
		})
	})

	When("with subnodes", func() {
		var subNode *Node
		BeforeEach(func() {
			subNode = NewNode(FunctionLevel, "sub")
			subNode.LinkParentNode(node)
		})
		It("should be cloned with subnodes", func() {
			checkTopNode(clonedNode)
			compareNode(clonedNode, node)
			Expect(clonedNode.SubNodes).To(HaveLen(1))
			compareNode(clonedNode.SubNodes[0], subNode)
		})
	})
})

func TestNode_Equal(t *testing.T) {
	node := NewNode(ModuleLevel, "top")

	type args struct {
		other *Node
	}
	tests := []struct {
		name  string
		nodeA *Node
		nodeB *Node
		want  bool
	}{
		{
			name:  "both nil node",
			nodeA: nil,
			nodeB: nil,
			want:  false,
		},
		{
			name:  "one of node is nil",
			nodeA: NewNode(ModuleLevel, "top"),
			nodeB: nil,
			want:  false,
		},
		{
			name:  "equal case 1",
			nodeA: node,
			nodeB: node,
			want:  true,
		},
		{
			name:  "equal case 2",
			nodeA: node,
			nodeB: NewNode(ModuleLevel, "top"),
			want:  true,
		},
		{
			name:  "not equal case",
			nodeA: node,
			nodeB: NewNode(ModuleLevel, "other"),
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.nodeA.Equal(tt.nodeB); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ = Describe("test for IdentifyLabels", func() {
	When("node have no parent node", func() {
		It("should return correct labels", func() {
			node := NewNode(ModuleLevel, "top")
			Expect(node.IdentifyLabels()).To(Equal(Labels{"module:top"}))
		})
	})

	When("node have parent node", func() {
		It("should return correct labels", func() {
			node := NewNode(ModuleLevel, "top")
			nodea := NewNode(ModuleLevel, "nodea")
			nodeb := NewNode(ModuleLevel, "nodeb")
			nodea.LinkParentNode(node)
			nodeb.LinkParentNode(nodea)
			Expect(nodea.IdentifyLabels()).To(Equal(Labels{"module:top#module:nodea"}))
			Expect(nodeb.IdentifyLabels()).To(Equal(Labels{"module:top#module:nodea#module:nodeb"}))
		})
	})
})
