package adapters

import "github.com/Sovianum/turbocycle/core/graph"

// NewTypedNode constructs Typed node out of its components
func NewTypedNode(node graph.Node, nodeType string) *TypedNode {
	return &TypedNode{
		NodeType: nodeType,
		Node:     node,
	}
}

// TypedNode is a helper struct combining Node with its type tag
type TypedNode struct {
	NodeType string
	Node     graph.Node
}
