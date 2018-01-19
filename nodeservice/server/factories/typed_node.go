package factories

import "github.com/Sovianum/turbocycle/core/graph"

func NewTypedNode(node graph.Node, nodeType string) *TypedNode {
	return &TypedNode{
		NodeType: nodeType,
		Node:     node,
	}
}

type TypedNode struct {
	NodeType string
	Node     graph.Node
}
