package adapters

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbocycle/core/graph"
)

type NodeAdapter interface {
	Create(data *pb.RequestData) (graph.Node, error)
	Update(node graph.Node, data *pb.RequestData) error
	GetState(node graph.Node, requiredFields []string) (*pb.NodeState, error)
	GetPort(tag string, node graph.Node) (graph.Port, error)
	GetPortDescriptions()[]PortDescription
}
