package mocks

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
)

type NodeAdapterMock struct {
	CreateFunc   func(data *pb.RequestData) (graph.Node, error)
	UpdateFunc   func(node graph.Node, data *pb.RequestData) error
	GetStateFunc func(node graph.Node, requiredFields []string) (*pb.NodeState, error)
	GetPortFunc  func(tag string, node graph.Node) (graph.Port, error)
	GetDecs      func() adapters.NodeDescription
}

func (m NodeAdapterMock) Create(data *pb.RequestData) (graph.Node, error) {
	return m.CreateFunc(data)
}

func (m NodeAdapterMock) Update(node graph.Node, data *pb.RequestData) error {
	return m.UpdateFunc(node, data)
}

func (m NodeAdapterMock) GetState(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
	return m.GetStateFunc(node, requiredFields)
}

func (m NodeAdapterMock) GetPort(tag string, node graph.Node) (graph.Port, error) {
	return m.GetPortFunc(tag, node)
}

func (m NodeAdapterMock) GetDescription() adapters.NodeDescription {
	return m.GetDecs()
}
