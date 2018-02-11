package mocks

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/pb"
)

// NodeAdapterMock mocks NodeAdapter interface
type NodeAdapterMock struct {
	CreateFunc   func(data *pb.RequestData) (graph.Node, error)
	UpdateFunc   func(node graph.Node, data *pb.RequestData) error
	GetStateFunc func(node graph.Node, requiredFields []string) (*pb.NodeState, error)
	GetPortFunc  func(tag string, node graph.Node) (graph.Port, error)
	GetDecs      func() *pb.NodeDescription
}

// Create mocks NodeAdapter.Create method
func (m NodeAdapterMock) Create(data *pb.RequestData) (graph.Node, error) {
	return m.CreateFunc(data)
}

// Update mocks NodeAdapter.Update method
func (m NodeAdapterMock) Update(node graph.Node, data *pb.RequestData) error {
	return m.UpdateFunc(node, data)
}

// GetState mocks NodeAdapter.GetState method
func (m NodeAdapterMock) GetState(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
	return m.GetStateFunc(node, requiredFields)
}

// GetPort mocks NodeAdapter.GetPort method
func (m NodeAdapterMock) GetPort(tag string, node graph.Node) (graph.Port, error) {
	return m.GetPortFunc(tag, node)
}

// GetDescription mocks NodeAdapter.GetDescription method
func (m NodeAdapterMock) GetDescription() *pb.NodeDescription {
	return m.GetDecs()
}
